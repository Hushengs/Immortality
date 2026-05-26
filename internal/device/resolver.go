package device

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type ResolveResult struct {
	ResolvedPath string
	Searched     []string
}

func ResolveADBPath(configured string) (ResolveResult, error) {
	configured = strings.TrimSpace(configured)
	if configured == "" {
		configured = "adb"
	}

	searched := []string{}
	try := func(candidate string) (string, bool) {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			return "", false
		}
		searched = append(searched, candidate)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, true
		}
		return "", false
	}

	if filepath.IsAbs(configured) {
		if resolved, ok := try(configured); ok {
			return ResolveResult{ResolvedPath: resolved, Searched: searched}, nil
		}
		return ResolveResult{Searched: searched}, fmt.Errorf("configured adb_path not found: %s", configured)
	}

	if lookedUp, err := exec.LookPath(configured); err == nil {
		searched = append(searched, "PATH:"+configured)
		return ResolveResult{ResolvedPath: lookedUp, Searched: searched}, nil
	}
	searched = append(searched, "PATH:"+configured)

	for _, candidate := range candidateADBPaths() {
		if resolved, ok := try(candidate); ok {
			return ResolveResult{ResolvedPath: resolved, Searched: searched}, nil
		}
	}

	return ResolveResult{Searched: searched}, fmt.Errorf("adb not found via PATH or common SDK locations")
}

func candidateADBPaths() []string {
	candidates := []string{}
	exeName := "adb"
	if runtime.GOOS == "windows" {
		exeName = "adb.exe"
	}

	appendSDK := func(base string) {
		base = strings.TrimSpace(base)
		if base == "" {
			return
		}
		candidates = append(candidates, filepath.Join(base, "platform-tools", exeName))
	}

	appendSDK(os.Getenv("ANDROID_HOME"))
	appendSDK(os.Getenv("ANDROID_SDK_ROOT"))
	appendSDK(filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk"))
	appendSDK(filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Android", "Sdk"))

	return uniqueStrings(candidates)
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
