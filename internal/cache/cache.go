package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/skymoore/vibe-zsh/internal/schema"
)

type Cache struct {
	dir string
	ttl time.Duration
}

type CacheEntry struct {
	Query     string                  `json:"query"`
	Response  *schema.CommandResponse `json:"response"`
	Timestamp time.Time               `json:"timestamp"`
}

func New(cacheDir string, ttl time.Duration) (*Cache, error) {
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

	return &Cache{
		dir: cacheDir,
		ttl: ttl,
	}, nil
}

func (c *Cache) Get(query string) (*schema.CommandResponse, bool) {
	key := c.hashQuery(query)
	path := filepath.Join(c.dir, key+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if time.Since(entry.Timestamp) > c.ttl {
		os.Remove(path)
		return nil, false
	}

	return entry.Response, true
}

func (c *Cache) Set(query string, response *schema.CommandResponse) error {
	key := c.hashQuery(query)
	path := filepath.Join(c.dir, key+".json")

	entry := CacheEntry{
		Query:     query,
		Response:  response,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Cache) hashQuery(query string) string {
	hash := sha256.Sum256([]byte(query))
	return hex.EncodeToString(hash[:])
}

func (c *Cache) Clear() error {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			os.Remove(filepath.Join(c.dir, entry.Name()))
		}
	}

	return nil
}
