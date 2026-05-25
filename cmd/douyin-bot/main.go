package main

import (
	"context"
	"flag"
	"log"
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

	configProvider := app.NewLocalConfigProvider(*configPath)
	config, err := configProvider.Load(ctx)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

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
