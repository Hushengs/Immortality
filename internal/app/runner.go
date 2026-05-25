package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Hushengs/Immortality/internal/actions"
	"github.com/Hushengs/Immortality/internal/capture"
	"github.com/Hushengs/Immortality/internal/detectors"
	"github.com/Hushengs/Immortality/internal/device"
	"github.com/Hushengs/Immortality/internal/ocr"
	"github.com/Hushengs/Immortality/internal/strategy"
)

type Runner struct {
	logger              *slog.Logger
	config              RuntimeConfig
	stateMachine        *StateMachine
	deviceManager       *device.Manager
	baseADB             *device.ADBClient
	ocrService          *ocr.Service
	fudaiDetector       *detectors.FudaiDetector
	popupLayoutDetector *detectors.PopupLayoutDetector
	rewardDetector      *detectors.RewardDetector
	liveStatusDetector  *detectors.LiveStatusDetector
	verifyDetector      *detectors.VerificationDetector
	prizeFilter         *strategy.PrizeFilter
	roomSwitchStrategy  *strategy.RoomSwitchStrategy
	recordRepo          RecordRepository
	deviceRepo          DeviceRepository
	liveRoomRepo        LiveRoomRepository
	management          ManagementReporter

	runCtx        RunContext
	adb           *device.ADBClient
	capture       *capture.ScreenshotService
	navigation    *actions.LiveNavigation
	lotteryAction *actions.LotteryActions
	rewardAction  *actions.RewardActions
	detection     FudaiDetection
	layout        PopupLayout
	analysis      PopupAnalysis
}

func NewRunner(
	logger *slog.Logger,
	config RuntimeConfig,
	baseADB *device.ADBClient,
	deviceManager *device.Manager,
	ocrService *ocr.Service,
	recordRepo RecordRepository,
	deviceRepo DeviceRepository,
	liveRoomRepo LiveRoomRepository,
	management ManagementReporter,
) *Runner {
	return &Runner{
		logger:              logger,
		config:              config,
		stateMachine:        NewStateMachine(),
		baseADB:             baseADB,
		deviceManager:       deviceManager,
		ocrService:          ocrService,
		fudaiDetector:       detectors.NewFudaiDetector(config, ocrService),
		popupLayoutDetector: detectors.NewPopupLayoutDetector(config),
		rewardDetector:      detectors.NewRewardDetector(config, ocrService),
		liveStatusDetector:  detectors.NewLiveStatusDetector(config, ocrService),
		verifyDetector:      detectors.NewVerificationDetector(config, ocrService),
		prizeFilter:         strategy.NewPrizeFilter(config.Strategy),
		roomSwitchStrategy:  strategy.NewRoomSwitchStrategy(config.MaxWaitMinutes),
		recordRepo:          recordRepo,
		deviceRepo:          deviceRepo,
		liveRoomRepo:        liveRoomRepo,
		management:          management,
		runCtx: RunContext{
			RunID: newID("run"),
		},
	}
}

func (r *Runner) Run(ctx context.Context) error {
	r.transition(StateSelectDevice)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if r.config.Debug.MaxIterations > 0 && r.runCtx.Iteration >= r.config.Debug.MaxIterations {
			r.logger.Info("max iterations reached, stop runner", "max_iterations", r.config.Debug.MaxIterations)
			return nil
		}

		switch r.stateMachine.Current() {
		case StateSelectDevice:
			if err := r.selectDevice(ctx); err != nil {
				return err
			}
			r.transition(StateCaptureScreen)
		case StateCaptureScreen:
			if err := r.captureScreen(ctx, "loop"); err != nil {
				return err
			}
			status, err := r.liveStatusDetector.Detect(ctx, r.runCtx.LastScreenshot)
			if err == nil && status.IsVerification {
				r.transition(StateHandleVerify)
				continue
			}
			if err == nil && status.IsLiveEnded {
				r.transition(StateRecoverState)
				continue
			}
			r.transition(StateCheckFudaiIcon)
		case StateCheckFudaiIcon:
			detection, err := r.fudaiDetector.Detect(ctx, r.runCtx.LastScreenshot)
			if err != nil {
				return err
			}
			r.detection = detection
			if !detection.Found {
				if r.runCtx.LastNoFudaiAt.IsZero() {
					r.runCtx.LastNoFudaiAt = time.Now()
				}
				if r.roomSwitchStrategy.ShouldSwitch(r.config, r.runCtx.LastNoFudaiAt) {
					r.transition(StateSwitchRoom)
				} else {
					if err := r.recordEvent(ctx, "fudai_not_found", detection.Reason, nil); err != nil {
						return err
					}
					if err := sleepWithContext(ctx, time.Duration(r.config.LoopIntervalSeconds)*time.Second); err != nil {
						return err
					}
					r.transition(StateCaptureScreen)
				}
				continue
			}
			r.runCtx.LastNoFudaiAt = time.Time{}
			r.transition(StateOpenFudaiDetail)
		case StateOpenFudaiDetail:
			if err := r.lotteryAction.OpenFudaiDetail(ctx, r.detection); err != nil {
				return err
			}
			if err := sleepWithContext(ctx, 1300*time.Millisecond); err != nil {
				return err
			}
			r.transition(StateParsePopupLayout)
		case StateParsePopupLayout:
			r.layout = r.popupLayoutDetector.Detect()
			r.transition(StateOCRContent)
		case StateOCRContent:
			if err := r.captureScreen(ctx, "popup"); err != nil {
				return err
			}
			analysis, err := r.rewardDetector.Analyze(ctx, r.runCtx.LastScreenshot, r.layout)
			if err != nil {
				return err
			}
			r.analysis = analysis
			r.runCtx.LastTask = LotteryTask{
				TaskID:          newID("task"),
				RunID:           r.runCtx.RunID,
				DeviceID:        r.runCtx.Device.DeviceID,
				RoomID:          r.runCtx.Room.RoomID,
				PrizeText:       analysis.PrizeText,
				CountdownText:   analysis.CountdownText,
				CountdownSecond: analysis.CountdownSeconds,
				ButtonText:      analysis.ButtonText,
				LayoutName:      analysis.LayoutName,
				DetectedAt:      time.Now(),
			}
			r.transition(StateDecideJoin)
		case StateDecideJoin:
			r.runCtx.LastDecision = r.prizeFilter.Decide(r.analysis)
			if !r.runCtx.LastDecision.ShouldJoin {
				if err := r.saveCurrentRecord(ctx, "skip", "pending"); err != nil {
					return err
				}
				_ = r.navigation.Back(ctx)
				if r.runCtx.LastDecision.ShouldSwitch {
					r.transition(StateSwitchRoom)
				} else {
					r.transition(StateCaptureScreen)
				}
				continue
			}
			r.transition(StateClickJoin)
		case StateClickJoin:
			if err := r.lotteryAction.ClickJoin(ctx, r.layout.ButtonRegion); err != nil {
				return err
			}
			r.runCtx.LastJoinAttemptAt = time.Now()
			r.transition(StateWaitResult)
		case StateWaitResult:
			waitSeconds := r.analysis.CountdownSeconds
			if waitSeconds < 3 {
				waitSeconds = 3
			}
			if waitSeconds > 15 {
				waitSeconds = 15
			}
			if err := sleepWithContext(ctx, time.Duration(waitSeconds)*time.Second); err != nil {
				return err
			}
			r.transition(StateCheckResult)
		case StateCheckResult:
			if err := r.captureScreen(ctx, "result"); err != nil {
				return err
			}
			result, err := r.rewardDetector.DetectResult(ctx, r.runCtx.LastScreenshot)
			if err != nil {
				return err
			}
			if err := r.saveCurrentRecord(ctx, "join", result); err != nil {
				return err
			}
			if result == "won" {
				r.transition(StateClaimReward)
			} else {
				_ = r.navigation.Back(ctx)
				r.transition(StateCaptureScreen)
			}
		case StateClaimReward:
			if err := r.rewardAction.Claim(ctx, r.layout.ButtonRegion); err != nil {
				return err
			}
			if err := r.recordEvent(ctx, "claim_reward", "reward flow reached", nil); err != nil {
				return err
			}
			_ = r.navigation.Back(ctx)
			r.transition(StateCaptureScreen)
		case StateSwitchRoom:
			if r.config.SwitchRoom {
				_ = r.navigation.Back(ctx)
				if err := sleepWithContext(ctx, 800*time.Millisecond); err != nil {
					return err
				}
				if err := r.navigation.SwitchRoom(ctx); err != nil {
					return err
				}
			}
			if err := sleepWithContext(ctx, 1200*time.Millisecond); err != nil {
				return err
			}
			r.transition(StateCaptureScreen)
		case StateRecoverState:
			if err := r.recordEvent(ctx, "recover_state", "try to recover from non-live page", nil); err != nil {
				return err
			}
			_ = r.navigation.Back(ctx)
			if err := sleepWithContext(ctx, 1200*time.Millisecond); err != nil {
				return err
			}
			r.transition(StateCaptureScreen)
		case StateHandleVerify:
			detected, text, err := r.verifyDetector.Detect(ctx, r.runCtx.LastScreenshot)
			if err != nil {
				return err
			}
			if detected {
				if err := r.recordEvent(ctx, "verification_detected", "manual intervention may be required", map[string]any{
					"ocr_text": text,
				}); err != nil {
					return err
				}
			}
			r.transition(StateRecoverState)
		default:
			return fmt.Errorf("unsupported state: %s", r.stateMachine.Current())
		}
	}
}

func (r *Runner) selectDevice(ctx context.Context) error {
	deviceInfo, adbWithSerial, err := r.deviceManager.SelectDevice(ctx, r.config.DeviceSerial)
	if err != nil {
		return err
	}

	r.adb = adbWithSerial
	r.capture = capture.NewScreenshotService(r.logger, adbWithSerial, r.config.Artifacts.ScreenshotDir, r.config.Debug)
	r.navigation = actions.NewLiveNavigation(r.logger, adbWithSerial, r.config)
	r.lotteryAction = actions.NewLotteryActions(r.logger, adbWithSerial)
	r.rewardAction = actions.NewRewardActions(r.logger, adbWithSerial, r.config)
	r.runCtx.Device = deviceInfo

	if r.config.Debug.DryRun {
		r.logger.Warn("runner started in dry-run mode", "static_screen_path", r.config.Debug.StaticScreenPath)
	}

	if err := r.deviceRepo.SaveHeartbeat(ctx, DeviceHeartbeat{
		RunID:      r.runCtx.RunID,
		DeviceID:   deviceInfo.DeviceID,
		Status:     "online",
		LastSeenAt: time.Now(),
	}); err != nil {
		return err
	}

	return r.management.ReportHeartbeat(ctx, DeviceHeartbeat{
		RunID:      r.runCtx.RunID,
		DeviceID:   deviceInfo.DeviceID,
		Status:     "online",
		LastSeenAt: time.Now(),
	})
}

func (r *Runner) captureScreen(ctx context.Context, prefix string) error {
	path, err := r.capture.Capture(ctx, prefix)
	if err != nil {
		return err
	}
	r.runCtx.LastScreenshot = path
	r.runCtx.Iteration++
	return nil
}

func (r *Runner) saveCurrentRecord(ctx context.Context, decision string, result string) error {
	record := LotteryRecord{
		RecordID:        newID("record"),
		RunID:           r.runCtx.RunID,
		TaskID:          r.runCtx.LastTask.TaskID,
		DeviceID:        r.runCtx.Device.DeviceID,
		RoomID:          r.runCtx.Room.RoomID,
		PrizeText:       r.analysis.PrizeText,
		CountdownText:   r.analysis.CountdownText,
		CountdownSecond: r.analysis.CountdownSeconds,
		ButtonText:      r.analysis.ButtonText,
		Decision:        decision,
		Result:          result,
		ScreenshotPath:  r.runCtx.LastScreenshot,
		Metadata: map[string]any{
			"reason": r.runCtx.LastDecision.Reason,
			"layout": r.analysis.LayoutName,
		},
		CreatedAt: time.Now(),
	}
	return r.recordRepo.SaveRecord(ctx, record)
}

func (r *Runner) recordEvent(ctx context.Context, eventType string, message string, fields map[string]any) error {
	event := RuntimeEvent{
		RunID:      r.runCtx.RunID,
		DeviceID:   r.runCtx.Device.DeviceID,
		RoomID:     r.runCtx.Room.RoomID,
		Type:       eventType,
		Message:    message,
		Screenshot: r.runCtx.LastScreenshot,
		Fields:     fields,
		CreatedAt:  time.Now(),
	}

	if err := r.recordRepo.SaveEvent(ctx, event); err != nil {
		return err
	}
	return r.management.ReportEvent(ctx, event)
}

func (r *Runner) transition(next State) {
	r.logger.Info("state transition", "from", r.stateMachine.Current(), "to", next)
	r.stateMachine.Transition(next)
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}
