package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hushengs/Immortality/internal/app"
	"github.com/Hushengs/Immortality/internal/device"
	"github.com/Hushengs/Immortality/internal/integration"
	"github.com/Hushengs/Immortality/internal/ocr"
	"github.com/Hushengs/Immortality/internal/repository"
	"github.com/Hushengs/Immortality/internal/utils"
)

func main() {
	configPath := flag.String("config", "./config.yaml", "path to config file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// #region agent log
	utils.WriteDebugLog("startup", "H2", "cmd/douyin-bot/main.go:24", "process startup", map[string]any{
		"configPath": *configPath,
		"pathEnv":    os.Getenv("PATH"),
	})
	// #endregion

	configProvider := app.NewLocalConfigProvider(*configPath)
	config, err := configProvider.Load(ctx)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	// #region agent log
	utils.WriteDebugLog("startup", "H1", "cmd/douyin-bot/main.go:36", "config loaded", map[string]any{
		"adbPath":         config.ADBPath,
		"dryRun":          config.Debug.DryRun,
		"staticScreen":    config.Debug.StaticScreenPath,
		"deviceSerial":    config.DeviceSerial,
		"commandTimeoutS": config.ADB.CommandTimeoutSeconds,
		"retryCount":      config.ADB.RetryCount,
	})
	// #endregion

	resolveResult, err := device.ResolveADBPath(config.ADBPath)
	// #region agent log
	utils.WriteDebugLog("startup", "H6", "cmd/douyin-bot/main.go:51", "adb path resolution attempted", map[string]any{
		"configuredADBPath": config.ADBPath,
		"resolvedADBPath":   resolveResult.ResolvedPath,
		"searched":          resolveResult.Searched,
		"resolveError":      errorString(err),
	})
	// #endregion
	if err != nil {
		log.Fatalf("resolve adb path failed: %v; searched=%v", err, resolveResult.Searched)
	}
	config.ADBPath = resolveResult.ResolvedPath

	if err := utils.EnsureArtifactDirs(config.Artifacts.LogDir, config.Artifacts.ScreenshotDir, config.Artifacts.RecordFile); err != nil {
		log.Fatalf("prepare artifact directories failed: %v", err)
	}

	logger, closeLogger, err := utils.NewLogger(config.Artifacts.LogDir)
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer func() {
		if closeLogger != nil {
			_ = closeLogger()
		}
	}()

	baseADB := device.NewADBClient(config.ADBPath, "", config.ADB, config.Debug.DryRun)
	// #region agent log
	utils.WriteDebugLog("startup", "H3", "cmd/douyin-bot/main.go:59", "adb client created", map[string]any{
		"adbPath": config.ADBPath,
		"dryRun":  config.Debug.DryRun,
	})
	// #endregion
	deviceManager := device.NewManager(logger, baseADB)
	ocrService := ocr.NewService(config.TesseractPath, config.OCR)
	recordRepo := repository.NewFileRecordRepository(config.Artifacts.RecordFile)
	deviceRepo := repository.NewMemoryDeviceRepository(logger)
	liveRoomRepo := repository.NewMemoryLiveRoomRepository(logger)
	managementClient := integration.NewNoopManagementClient(logger)

	runner := app.NewRunner(
		logger,
		config,
		baseADB,
		deviceManager,
		ocrService,
		recordRepo,
		deviceRepo,
		liveRoomRepo,
		managementClient,
	)

	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		logger.Error("runner exited with error", "error", err)
		log.Fatal(err)
	}
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
