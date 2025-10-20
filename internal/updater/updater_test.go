package updater

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStateManagement(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	state, err := LoadState()
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	if !state.CheckEnabled {
		t.Error("CheckEnabled should be true by default")
	}

	state.CurrentVersion = "v1.0.0"
	state.AvailableVersion = "v1.0.1"
	state.LastCheck = time.Now()

	if err := SaveState(state); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	loaded, err := LoadState()
	if err != nil {
		t.Fatalf("LoadState after save failed: %v", err)
	}

	if loaded.CurrentVersion != "v1.0.0" {
		t.Errorf("Expected CurrentVersion v1.0.0, got %s", loaded.CurrentVersion)
	}

	if loaded.AvailableVersion != "v1.0.1" {
		t.Errorf("Expected AvailableVersion v1.0.1, got %s", loaded.AvailableVersion)
	}

	configPath := filepath.Join(tmpDir, ".config", "vibe", "state.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("State file was not created")
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		current  string
		latest   string
		expected bool
	}{
		{"v1.0.0", "v1.0.1", true},
		{"v1.0.1", "v1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"1.0.0", "1.0.1", true},
		{"v1.2.3", "v2.0.0", true},
		{"dev", "v1.0.0", false},
		{"v0.1.0", "v0.1.1", true},
	}

	for _, tt := range tests {
		result := compareVersions(tt.current, tt.latest)
		if result != tt.expected {
			t.Errorf("compareVersions(%s, %s) = %v, expected %v",
				tt.current, tt.latest, result, tt.expected)
		}
	}
}

func TestShouldCheckForUpdates(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	state := &State{
		CheckEnabled: true,
		LastCheck:    time.Now().Add(-8 * 24 * time.Hour),
	}

	if !shouldCheckForUpdates(state) {
		t.Error("Should check for updates after 8 days")
	}

	state.LastCheck = time.Now()
	if shouldCheckForUpdates(state) {
		t.Error("Should not check for updates immediately after last check")
	}

	state.CheckEnabled = false
	state.LastCheck = time.Now().Add(-8 * 24 * time.Hour)
	if shouldCheckForUpdates(state) {
		t.Error("Should not check when disabled")
	}

	state.CheckEnabled = true
	os.Setenv("VIBE_AUTO_UPDATE", "false")
	defer os.Unsetenv("VIBE_AUTO_UPDATE")
	if shouldCheckForUpdates(state) {
		t.Error("Should not check when VIBE_AUTO_UPDATE is false")
	}
}

func TestGetBinaryName(t *testing.T) {
	name := getBinaryName()
	if name == "" {
		t.Error("Binary name should not be empty")
	}

	if len(name) < 5 {
		t.Errorf("Binary name seems too short: %s", name)
	}
}

func TestGetDownloadURL(t *testing.T) {
	url := getDownloadURL("v1.2.3")
	expected := "https://github.com/skymoore/vibe-zsh/releases/download/v1.2.3/vibe-"

	if len(url) < len(expected) {
		t.Errorf("Download URL seems too short: %s", url)
	}

	if url[:len(expected)] != expected {
		t.Errorf("Download URL has wrong prefix: %s", url)
	}
}

func TestGetInstallMethod(t *testing.T) {
	method := getInstallMethod()

	validMethods := []string{"oh-my-zsh", "standalone", "unknown"}
	valid := false
	for _, m := range validMethods {
		if method == m {
			valid = true
			break
		}
	}

	if !valid {
		t.Errorf("Invalid install method: %s", method)
	}
}

func TestComputeChecksum(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test.txt"

	content := []byte("hello world")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	checksum, err := computeChecksum(testFile)
	if err != nil {
		t.Fatalf("computeChecksum failed: %v", err)
	}

	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if checksum != expected {
		t.Errorf("Expected checksum %s, got %s", expected, checksum)
	}
}

func TestVerifyChecksum(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test.txt"

	content := []byte("hello world")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	correctChecksum := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if err := verifyChecksum(testFile, correctChecksum); err != nil {
		t.Errorf("verifyChecksum should pass with correct checksum: %v", err)
	}

	wrongChecksum := "0000000000000000000000000000000000000000000000000000000000000000"
	if err := verifyChecksum(testFile, wrongChecksum); err == nil {
		t.Error("verifyChecksum should fail with wrong checksum")
	}
}

func TestGetChecksumsURL(t *testing.T) {
	url := getChecksumsURL("v1.2.3")
	expected := "https://github.com/skymoore/vibe-zsh/releases/download/v1.2.3/checksums.txt"

	if url != expected {
		t.Errorf("Expected URL %s, got %s", expected, url)
	}
}
