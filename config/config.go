package config

import (
	"os"
	"time"
)

type Config struct {
	Port                     string
	CacheRefreshInterval     string
	CacheTTL                 time.Duration
	FlashscoreABAFixturesURL string
	FlashscoreABAResultsURL  string
	EuroleagueAPIURL         string
	Team                     string
	TeamCodeEuroleague       string
	DefaultGameDuration      int
}

func Load() *Config {
	port := getEnv("PORT", "3000")
	cacheRefreshInterval := getEnv("CACHE_REFRESH_INTERVAL", "0 0 */2 * *")
	cacheTTL := 2 * 24 * time.Hour

	return &Config{
		Port:                     port,
		CacheRefreshInterval:     cacheRefreshInterval,
		CacheTTL:                 cacheTTL,
		FlashscoreABAFixturesURL: "https://www.flashscore.com/basketball/europe/aba-league/fixtures/",
		FlashscoreABAResultsURL:  "https://www.flashscore.com/basketball/europe/aba-league/results/",
		EuroleagueAPIURL:         "https://feeds.incrowdsports.com/provider/euroleague-feeds/v2/competitions/E/seasons/E2025/games",
		Team:                     "Partizan",
		TeamCodeEuroleague:       "PAR",
		DefaultGameDuration:      2,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
