package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()

	h, err := New(tmpDir, 100)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if h.dir != tmpDir {
		t.Errorf("dir = %v, want %v", h.dir, tmpDir)
	}

	if h.maxSize != 100 {
		t.Errorf("maxSize = %v, want 100", h.maxSize)
	}

	// Check that directory was created
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Errorf("directory was not created")
	}
}

func TestAddAndList(t *testing.T) {
	tmpDir := t.TempDir()
	h, err := New(tmpDir, 100)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Add first entry
	err = h.Add("list files", "ls -la")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// List entries
	entries, err := h.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("len(entries) = %v, want 1", len(entries))
	}

	if entries[0].Query != "list files" {
		t.Errorf("Query = %v, want 'list files'", entries[0].Query)
	}

	if entries[0].Command != "ls -la" {
		t.Errorf("Command = %v, want 'ls -la'", entries[0].Command)
	}

	// Add second entry
	err = h.Add("show processes", "ps aux")
	if err != nil {
		t.Fatalf("Add() second error = %v", err)
	}

	entries, err = h.List()
	if err != nil {
		t.Fatalf("List() second error = %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("len(entries) = %v, want 2", len(entries))
	}

	// Most recent should be first
	if entries[0].Query != "show processes" {
		t.Errorf("First entry Query = %v, want 'show processes'", entries[0].Query)
	}

	if entries[1].Query != "list files" {
		t.Errorf("Second entry Query = %v, want 'list files'", entries[1].Query)
	}
}

func TestMaxSize(t *testing.T) {
	tmpDir := t.TempDir()
	maxSize := 5
	h, err := New(tmpDir, maxSize)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Add more entries than maxSize
	for i := 0; i < 10; i++ {
		err = h.Add("query", "command")
		if err != nil {
			t.Fatalf("Add() error = %v", err)
		}
	}

	entries, err := h.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(entries) != maxSize {
		t.Errorf("len(entries) = %v, want %v", len(entries), maxSize)
	}
}

func TestClear(t *testing.T) {
	tmpDir := t.TempDir()
	h, err := New(tmpDir, 100)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Add entries
	err = h.Add("query1", "command1")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	err = h.Add("query2", "command2")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Clear
	err = h.Clear()
	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	// Check that history is empty
	entries, err := h.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("len(entries) = %v, want 0", len(entries))
	}
}

func TestListEmptyHistory(t *testing.T) {
	tmpDir := t.TempDir()
	h, err := New(tmpDir, 100)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	entries, err := h.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("len(entries) = %v, want 0", len(entries))
	}
}

func TestCorruptedHistoryFile(t *testing.T) {
	tmpDir := t.TempDir()
	h, err := New(tmpDir, 100)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Write corrupted JSON to history file
	corruptedData := []byte("this is not valid json")
	err = os.WriteFile(filepath.Join(tmpDir, "history.json"), corruptedData, 0644)
	if err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// List should return empty array, not error
	entries, err := h.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("len(entries) = %v, want 0", len(entries))
	}

	// Adding new entry should work and fix the file
	err = h.Add("new query", "new command")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	entries, err = h.List()
	if err != nil {
		t.Fatalf("List() after Add error = %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("len(entries) = %v, want 1", len(entries))
	}
}

func TestTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	h, err := New(tmpDir, 100)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	before := time.Now()
	err = h.Add("query", "command")
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	after := time.Now()

	entries, err := h.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("len(entries) = %v, want 1", len(entries))
	}

	timestamp := entries[0].Timestamp
	if timestamp.Before(before) || timestamp.After(after) {
		t.Errorf("Timestamp %v is not between %v and %v", timestamp, before, after)
	}
}
