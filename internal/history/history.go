package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	Query     string    `json:"query"`
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

type History struct {
	dir      string
	maxSize  int
	filePath string
}

func New(cacheDir string, maxSize int) (*History, error) {
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cacheDir = filepath.Join(home, ".cache", "vibe")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	return &History{
		dir:      cacheDir,
		maxSize:  maxSize,
		filePath: filepath.Join(cacheDir, "history.json"),
	}, nil
}

func (h *History) Add(query, command string) error {
	entries, err := h.List()
	if err != nil {
		// If file doesn't exist or is corrupted, start fresh
		entries = []Entry{}
	}

	// Add new entry at the beginning (most recent first)
	newEntry := Entry{
		Query:     query,
		Command:   command,
		Timestamp: time.Now(),
		Count:     1,
	}

	entries = append([]Entry{newEntry}, entries...)

	// Enforce max size limit
	if len(entries) > h.maxSize {
		entries = entries[:h.maxSize]
	}

	return h.save(entries)
}

func (h *History) List() ([]Entry, error) {
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		// If corrupted, return empty and let caller handle
		return []Entry{}, nil
	}

	return entries, nil
}

func (h *History) Clear() error {
	return h.save([]Entry{})
}

func (h *History) GetDir() string {
	return h.dir
}

func (h *History) save(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	// Atomic write: write to temp file, then rename
	tmpPath := h.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, h.filePath)
}
