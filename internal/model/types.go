package model

import (
	"context"
	"time"
)

type Point struct {
	X int `json:"x" yaml:"x"`
	Y int `json:"y" yaml:"y"`
}

type Region struct {
	X      int `json:"x" yaml:"x"`
	Y      int `json:"y" yaml:"y"`
	Width  int `json:"width" yaml:"width"`
	Height int `json:"height" yaml:"height"`
}

func (r Region) Center() Point {
	return Point{
		X: r.X + (r.Width / 2),
		Y: r.Y + (r.Height / 2),
	}
}

type Resolution struct {
	Width  int `json:"width" yaml:"width"`
	Height int `json:"height" yaml:"height"`
}

type ArtifactConfig struct {
	ScreenshotDir string `json:"screenshot_dir" yaml:"screenshot_dir"`
	LogDir        string `json:"log_dir" yaml:"log_dir"`
	RecordFile    string `json:"record_file" yaml:"record_file"`
}

type ADBConfig struct {
	CommandTimeoutSeconds int `json:"command_timeout_seconds" yaml:"command_timeout_seconds"`
	RetryCount            int `json:"retry_count" yaml:"retry_count"`
	RetryDelayMS          int `json:"retry_delay_ms" yaml:"retry_delay_ms"`
}

type DebugConfig struct {
	DryRun           bool   `json:"dry_run" yaml:"dry_run"`
	MaxIterations    int    `json:"max_iterations" yaml:"max_iterations"`
	StaticScreenPath string `json:"static_screen_path" yaml:"static_screen_path"`
}

type OCRConfig struct {
	DefaultLang       string `json:"default_lang" yaml:"default_lang"`
	DigitsLang        string `json:"digits_lang" yaml:"digits_lang"`
	PSM               int    `json:"psm" yaml:"psm"`
	PreprocessEnabled bool   `json:"preprocess_enabled" yaml:"preprocess_enabled"`
	ThresholdValue    uint8  `json:"threshold_value" yaml:"threshold_value"`
}

type Thresholds struct {
	FudaiRedRatio       float64 `json:"fudai_red_ratio" yaml:"fudai_red_ratio"`
	TextMatchConfidence float64 `json:"text_match_confidence" yaml:"text_match_confidence"`
}

type StrategyConfig struct {
	PrizeWhitelist []string `json:"prize_whitelist" yaml:"prize_whitelist"`
	PrizeBlacklist []string `json:"prize_blacklist" yaml:"prize_blacklist"`
	MinJoinSeconds int      `json:"min_join_seconds" yaml:"min_join_seconds"`
	MaxJoinSeconds int      `json:"max_join_seconds" yaml:"max_join_seconds"`
	JoinKeywords   []string `json:"join_keywords" yaml:"join_keywords"`
}

type ScreenRegionSet struct {
	FudaiIcon      Region `json:"fudai_icon" yaml:"fudai_icon"`
	PopupPrize     Region `json:"popup_prize" yaml:"popup_prize"`
	PopupCountdown Region `json:"popup_countdown" yaml:"popup_countdown"`
	PopupButton    Region `json:"popup_button" yaml:"popup_button"`
	LiveStatus     Region `json:"live_status" yaml:"live_status"`
	Verification   Region `json:"verification" yaml:"verification"`
	SwitchSwipe    Region `json:"switch_swipe" yaml:"switch_swipe"`
	FirstRoomTap   Region `json:"first_room_tap" yaml:"first_room_tap"`
}

type KeywordSet struct {
	Fudai        []string `json:"fudai" yaml:"fudai"`
	LiveEnded    []string `json:"live_ended" yaml:"live_ended"`
	Verification []string `json:"verification" yaml:"verification"`
	Won          []string `json:"won" yaml:"won"`
}

type RuntimeConfig struct {
	AppName             string          `json:"app_name" yaml:"app_name"`
	SafeMode            bool            `json:"safe_mode" yaml:"safe_mode"`
	AutoClaimReward     bool            `json:"auto_claim_reward" yaml:"auto_claim_reward"`
	AutoPlaceOrder      bool            `json:"auto_place_order" yaml:"auto_place_order"`
	SwitchRoom          bool            `json:"switch_room" yaml:"switch_room"`
	MaxWaitMinutes      int             `json:"max_wait_minutes" yaml:"max_wait_minutes"`
	LoopIntervalSeconds int             `json:"loop_interval_seconds" yaml:"loop_interval_seconds"`
	DeviceSerial        string          `json:"device_serial" yaml:"device_serial"`
	ADBPath             string          `json:"adb_path" yaml:"adb_path"`
	ADB                 ADBConfig       `json:"adb" yaml:"adb"`
	TesseractPath       string          `json:"tesseract_path" yaml:"tesseract_path"`
	Debug               DebugConfig     `json:"debug" yaml:"debug"`
	Resolution          Resolution      `json:"resolution" yaml:"resolution"`
	YOffset             int             `json:"y_offset" yaml:"y_offset"`
	Artifacts           ArtifactConfig  `json:"artifacts" yaml:"artifacts"`
	Strategy            StrategyConfig  `json:"strategy" yaml:"strategy"`
	OCR                 OCRConfig       `json:"ocr" yaml:"ocr"`
	ScreenRegions       ScreenRegionSet `json:"screen_regions" yaml:"screen_regions"`
	Thresholds          Thresholds      `json:"thresholds" yaml:"thresholds"`
	Keywords            KeywordSet      `json:"keywords" yaml:"keywords"`
}

type Device struct {
	DeviceID   string    `json:"device_id"`
	Serial     string    `json:"serial"`
	Model      string    `json:"model"`
	Online     bool      `json:"online"`
	LastSeenAt time.Time `json:"last_seen_at"`
}

type LiveRoom struct {
	RoomID        string    `json:"room_id"`
	Title         string    `json:"title"`
	Anchor        string    `json:"anchor"`
	Status        string    `json:"status"`
	LastActiveAt  time.Time `json:"last_active_at"`
	LastFudaiSeen time.Time `json:"last_fudai_seen"`
}

type LotteryTask struct {
	TaskID          string    `json:"task_id"`
	RunID           string    `json:"run_id"`
	DeviceID        string    `json:"device_id"`
	RoomID          string    `json:"room_id"`
	PrizeText       string    `json:"prize_text"`
	CountdownText   string    `json:"countdown_text"`
	CountdownSecond int       `json:"countdown_seconds"`
	ButtonText      string    `json:"button_text"`
	LayoutName      string    `json:"layout_name"`
	DetectedAt      time.Time `json:"detected_at"`
}

type LotteryRecord struct {
	RecordID        string         `json:"record_id"`
	RunID           string         `json:"run_id"`
	TaskID          string         `json:"task_id"`
	DeviceID        string         `json:"device_id"`
	RoomID          string         `json:"room_id"`
	PrizeText       string         `json:"prize_text"`
	CountdownText   string         `json:"countdown_text"`
	CountdownSecond int            `json:"countdown_seconds"`
	ButtonText      string         `json:"button_text"`
	Decision        string         `json:"decision"`
	Result          string         `json:"result"`
	ScreenshotPath  string         `json:"screenshot_path"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
}

type RuntimeEvent struct {
	RunID      string         `json:"run_id"`
	DeviceID   string         `json:"device_id"`
	RoomID     string         `json:"room_id,omitempty"`
	Type       string         `json:"type"`
	Message    string         `json:"message"`
	Screenshot string         `json:"screenshot,omitempty"`
	Fields     map[string]any `json:"fields,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
}

type DeviceHeartbeat struct {
	RunID      string    `json:"run_id"`
	DeviceID   string    `json:"device_id"`
	Status     string    `json:"status"`
	LastSeenAt time.Time `json:"last_seen_at"`
}

type FudaiDetection struct {
	Found        bool
	MatchScore   float64
	Reason       string
	ClickPoint   Point
	SampleRegion Region
}

type PopupLayout struct {
	Name            string
	PrizeRegion     Region
	CountdownRegion Region
	ButtonRegion    Region
}

type PopupAnalysis struct {
	PrizeText        string
	CountdownText    string
	CountdownSeconds int
	ButtonText       string
	NeedTask         bool
	NeedFansClub     bool
	CannotJoin       bool
	AlreadyJoined    bool
	LayoutName       string
}

type Decision struct {
	ShouldJoin   bool
	ShouldClaim  bool
	ShouldSkip   bool
	ShouldSwitch bool
	Reason       string
}

type PageStatus struct {
	Kind           string
	IsLiveRoom     bool
	IsLiveEnded    bool
	IsVerification bool
	IsListPage     bool
	IsFollowPage   bool
}

type RunContext struct {
	RunID             string
	Device            Device
	Room              LiveRoom
	LastScreenshot    string
	LastTask          LotteryTask
	LastDecision      Decision
	LastNoFudaiAt     time.Time
	LastJoinAttemptAt time.Time
	Iteration         int
}

type ConfigProvider interface {
	Load(context.Context) (RuntimeConfig, error)
}

type RecordRepository interface {
	SaveRecord(context.Context, LotteryRecord) error
	SaveEvent(context.Context, RuntimeEvent) error
}

type DeviceRepository interface {
	SaveHeartbeat(context.Context, DeviceHeartbeat) error
}

type LiveRoomRepository interface {
	Save(context.Context, LiveRoom) error
}

type ManagementReporter interface {
	ReportEvent(context.Context, RuntimeEvent) error
	ReportHeartbeat(context.Context, DeviceHeartbeat) error
}
