package updater

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func getArchiveName(version string) string {
	return fmt.Sprintf("vibe-zsh-%s-%s-%s.tar.gz", version, runtime.GOOS, runtime.GOARCH)
}

func getBinaryName() string {
	return "vibe"
}

func getDownloadURL(version string) string {
	archiveName := getArchiveName(version)
	return fmt.Sprintf("https://github.com/skymoore/vibe-zsh/releases/download/%s/%s", version, archiveName)
}

func getChecksumsURL(version string) string {
	return fmt.Sprintf("https://github.com/skymoore/vibe-zsh/releases/download/%s/checksums.txt", version)
}

func getCurrentBinaryPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(execPath)
}

func downloadFile(url string, dest string) error {
	client := &http.Client{
		Timeout: 2 * time.Minute,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, destPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Name == "vibe" || strings.HasSuffix(header.Name, "/vibe") {
			out, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer out.Close()

			if _, err := io.Copy(out, tr); err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("vibe binary not found in archive")
}

func verifyBinary(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.Size() < 1024 {
		return fmt.Errorf("downloaded binary is too small (%d bytes)", info.Size())
	}

	return nil
}

func backupBinary(path string) (string, error) {
	backupPath := path + ".backup"

	input, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(backupPath, input, 0755)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

func PerformUpdate(currentVersion string) error {
	state, err := LoadState()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if state.AvailableVersion == "" {
		fmt.Println("Checking for updates...")
		latestVersion, err := checkLatestVersion()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		if !compareVersions(currentVersion, latestVersion) {
			fmt.Printf("Already on latest version %s\n", currentVersion)
			return nil
		}

		state.AvailableVersion = latestVersion
	}

	targetVersion := state.AvailableVersion
	fmt.Printf("Updating from %s to %s...\n", currentVersion, targetVersion)

	currentBinary, err := getCurrentBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to get current binary path: %w", err)
	}

	fmt.Println("Downloading checksums...")
	checksums, err := downloadChecksums(targetVersion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not download checksums: %v\n", err)
		fmt.Println("Continuing without checksum verification...")
		checksums = nil
	}

	tmpArchive := filepath.Join(os.TempDir(), fmt.Sprintf("vibe-archive-%d.tar.gz", time.Now().Unix()))
	defer os.Remove(tmpArchive)

	downloadURL := getDownloadURL(targetVersion)
	fmt.Printf("Downloading %s...\n", downloadURL)

	if err := downloadFile(downloadURL, tmpArchive); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if checksums != nil {
		archiveName := getArchiveName(targetVersion)
		if expectedChecksum, ok := checksums[archiveName]; ok {
			fmt.Println("Verifying checksum...")
			if err := verifyChecksum(tmpArchive, expectedChecksum); err != nil {
				return fmt.Errorf("checksum verification failed: %w", err)
			}
			fmt.Println("✓ Checksum verified")
		} else {
			fmt.Fprintf(os.Stderr, "Warning: No checksum found for %s\n", archiveName)
		}
	}

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("vibe-new-%d", time.Now().Unix()))
	defer os.Remove(tmpFile)

	fmt.Println("Extracting binary...")
	if err := extractBinary(tmpArchive, tmpFile); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	if err := verifyBinary(tmpFile); err != nil {
		return fmt.Errorf("binary verification failed: %w", err)
	}

	if err := os.Chmod(tmpFile, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	fmt.Println("Creating backup...")
	backupPath, err := backupBinary(currentBinary)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Println("Installing new binary...")
	if err := os.Rename(tmpFile, currentBinary); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	state.CurrentVersion = targetVersion
	state.AvailableVersion = ""
	state.NotificationShown = false
	state.LastCheck = time.Now()

	if err := SaveState(state); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to update state: %v\n", err)
	}

	fmt.Printf("✅ Successfully updated to %s\n", targetVersion)
	fmt.Printf("Backup saved to: %s\n", backupPath)

	return nil
}

func computeChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func downloadChecksums(version string) (map[string]string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url := getChecksumsURL(version)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download checksums: status %d", resp.StatusCode)
	}

	checksums := make(map[string]string)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			checksums[parts[1]] = parts[0]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return checksums, nil
}

func verifyChecksum(filePath string, expectedChecksum string) error {
	actualChecksum, err := computeChecksum(filePath)
	if err != nil {
		return fmt.Errorf("failed to compute checksum: %w", err)
	}

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

func getInstallMethod() string {
	execPath, err := os.Executable()
	if err != nil {
		return "unknown"
	}

	path, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return "unknown"
	}

	if strings.Contains(path, ".oh-my-zsh/custom/plugins/vibe") {
		return "oh-my-zsh"
	}

	return "standalone"
}
