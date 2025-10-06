package cache

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/milicacurcic/bcpartizan-ics/types"
)

const cacheFile = "data/cache.json"

type CacheData struct {
	Timestamp int64        `json:"timestamp"`
	Games     []types.Game `json:"games"`
}

func Load(cacheTTL time.Duration) ([]types.Game, error) {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	if cache.Timestamp > 0 && (now-cache.Timestamp) < int64(cacheTTL.Seconds()) {
		return cache.Games, nil
	}

	return nil, nil
}

func Save(games []types.Game) error {
	cache := CacheData{
		Timestamp: time.Now().Unix(),
		Games:     games,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(cacheFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating cache directory: %v", err)
		return err
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		log.Printf("Error saving cache: %v", err)
		return err
	}

	return nil
}
