package device

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Hushengs/Immortality/internal/model"
)

type ADBClient struct {
	adbPath string
	serial  string
	config  model.ADBConfig
	dryRun  bool
}

func NewADBClient(adbPath, serial string, config model.ADBConfig, dryRun bool) *ADBClient {
	return &ADBClient{
		adbPath: adbPath,
		serial:  serial,
		config:  normalizeADBConfig(config),
		dryRun:  dryRun,
	}
}

func (c *ADBClient) WithSerial(serial string) *ADBClient {
	return &ADBClient{adbPath: c.adbPath, serial: serial, config: c.config, dryRun: c.dryRun}
}

func (c *ADBClient) ListDevices(ctx context.Context) ([]model.Device, error) {
	if c.dryRun {
		return []model.Device{{
			DeviceID: "dry-run-device",
			Serial:   "dry-run-device",
			Online:   true,
		}}, nil
	}

	output, err := c.runADBCommand(ctx, "devices")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	devices := make([]model.Device, 0, len(lines))
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		devices = append(devices, model.Device{
			DeviceID: fields[0],
			Serial:   fields[0],
			Online:   fields[1] == "device",
		})
	}

	return devices, nil
}

func (c *ADBClient) Tap(ctx context.Context, point model.Point) error {
	_, err := c.Shell(ctx, "input", "tap", fmt.Sprintf("%d", point.X), fmt.Sprintf("%d", point.Y))
	return err
}

func (c *ADBClient) Swipe(ctx context.Context, from, to model.Point, durationMS int) error {
	_, err := c.Shell(
		ctx,
		"input", "swipe",
		fmt.Sprintf("%d", from.X), fmt.Sprintf("%d", from.Y),
		fmt.Sprintf("%d", to.X), fmt.Sprintf("%d", to.Y),
		fmt.Sprintf("%d", durationMS),
	)
	return err
}

func (c *ADBClient) KeyEvent(ctx context.Context, key string) error {
	_, err := c.Shell(ctx, "input", "keyevent", key)
	return err
}

func (c *ADBClient) Shell(ctx context.Context, args ...string) (string, error) {
	if c.dryRun {
		return "dry-run shell skipped", nil
	}

	commandArgs := c.baseArgs()
	commandArgs = append(commandArgs, "shell")
	commandArgs = append(commandArgs, args...)

	output, err := c.runADBCommand(ctx, commandArgs...)
	return strings.TrimSpace(string(output)), err
}

func (c *ADBClient) Screencap(ctx context.Context, outputPath string) error {
	if c.dryRun {
		return errors.New("dry-run screencap requires a static screenshot source")
	}

	commandArgs := c.baseArgs()
	commandArgs = append(commandArgs, "exec-out", "screencap", "-p")

	raw, err := c.runADBCommand(ctx, commandArgs...)
	if err != nil {
		return err
	}

	raw = bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n"))
	return os.WriteFile(outputPath, raw, 0o644)
}

func (c *ADBClient) baseArgs() []string {
	if c.serial == "" {
		return nil
	}
	return []string{"-s", c.serial}
}

func (c *ADBClient) runADBCommand(ctx context.Context, args ...string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.RetryCount; attempt++ {
		commandCtx := ctx
		cancel := func() {}
		if c.config.CommandTimeoutSeconds > 0 {
			commandCtx, cancel = context.WithTimeout(ctx, time.Duration(c.config.CommandTimeoutSeconds)*time.Second)
		}

		cmd := exec.CommandContext(commandCtx, c.adbPath, args...)
		output, err := cmd.CombinedOutput()
		cancel()
		if err == nil {
			return output, nil
		}

		lastErr = fmt.Errorf("adb %s failed: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
		if attempt == c.config.RetryCount {
			break
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(c.config.RetryDelayMS) * time.Millisecond):
		}
	}

	return nil, lastErr
}

func normalizeADBConfig(config model.ADBConfig) model.ADBConfig {
	if config.CommandTimeoutSeconds <= 0 {
		config.CommandTimeoutSeconds = 15
	}
	if config.RetryCount < 0 {
		config.RetryCount = 0
	}
	if config.RetryDelayMS <= 0 {
		config.RetryDelayMS = 800
	}
	return config
}
