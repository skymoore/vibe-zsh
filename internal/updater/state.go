package updater

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	LastCheck         time.Time `json:"last_check"`
	CurrentVersion    string    `json:"current_version"`
	AvailableVersion  string    `json:"available_version"`
	NotificationShown bool      `json:"notification_shown"`
	CheckEnabled      bool      `json:"check_enabled"`
}

func getStateFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "vibe")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "state.json"), nil
}

func LoadState() (*State, error) {
	path, err := getStateFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{
				CheckEnabled: true,
			}, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return &State{
			CheckEnabled: true,
		}, nil
	}

	return &state, nil
}

func SaveState(state *State) error {
	path, err := getStateFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
