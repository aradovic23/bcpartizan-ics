package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/milicacurcic/bcpartizan-ics/cache"
	"github.com/milicacurcic/bcpartizan-ics/config"
	"github.com/milicacurcic/bcpartizan-ics/ics"
	"github.com/milicacurcic/bcpartizan-ics/scraper"
	"github.com/milicacurcic/bcpartizan-ics/types"
	"github.com/robfig/cron/v3"
)

var (
	cachedGames []types.Game
	cacheMutex  sync.RWMutex
	cfg         *config.Config
)

func refreshSchedule() []types.Game {
	log.Println("Refreshing schedule...")
	games := scraper.FetchAllSchedules(cfg)

	if len(games) == 0 {
		log.Println("No games fetched from real sources, using mock data")
		games = scraper.GetMockSchedule()
	}

	cacheMutex.Lock()
	cachedGames = games
	cacheMutex.Unlock()

	if err := cache.Save(games); err != nil {
		log.Printf("Error saving cache: %v", err)
	}

	log.Printf("Schedule refreshed: %d games found", len(games))
	return games
}

func getGames() []types.Game {
	cacheMutex.RLock()
	if len(cachedGames) > 0 {
		games := cachedGames
		cacheMutex.RUnlock()
		return games
	}
	cacheMutex.RUnlock()

	cached, err := cache.Load(cfg.CacheTTL)
	if err == nil && cached != nil {
		cacheMutex.Lock()
		cachedGames = cached
		cacheMutex.Unlock()
		return cached
	}

	return refreshSchedule()
}

func calendarHandler(w http.ResponseWriter, r *http.Request) {
	games := getGames()
	icsContent, err := ics.GenerateCalendar(games, cfg)
	if err != nil {
		log.Printf("Error generating calendar: %v", err)
		http.Error(w, "Error generating calendar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Cache-Control", "max-age=3600")
	w.Write([]byte(icsContent))
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {
	games := getGames()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(games),
		"games": games,
	})
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	games := refreshSchedule()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Schedule refreshed",
		"count":   len(games),
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	cacheMutex.RLock()
	gameCount := len(cachedGames)
	cacheMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "healthy",
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"games_cached": gameCount,
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	html := fmt.Sprintf(`
<html>
  <head><title>Partizan Basketball Schedule</title></head>
  <body style="font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px;">
    <h1>🏀 KK Partizan Schedule Calendar</h1>
    <p>Subscribe to Partizan basketball games across all competitions.</p>
    
    <h2>Calendar Subscription URL:</h2>
    <code style="background: #f4f4f4; padding: 10px; display: block; margin: 10px 0;">
      %s://%s/calendar.ics
    </code>
    
    <h3>How to subscribe:</h3>
    <ul>
      <li><strong>Apple Calendar:</strong> File → New Calendar Subscription → Paste URL</li>
      <li><strong>Google Calendar:</strong> Settings → Add Calendar → From URL → Paste URL</li>
      <li><strong>Outlook:</strong> Add Calendar → Subscribe from web → Paste URL</li>
    </ul>
    
    <h3>Features:</h3>
    <ul>
      <li>✅ All Partizan games across Euroleague, ABA League, and domestic competitions</li>
      <li>✅ Game title includes competition name and teams</li>
      <li>✅ Venue location information</li>
      <li>✅ 30-minute reminder before each game</li>
      <li>✅ Automatically updates every 6 hours</li>
    </ul>
    
    <p><a href="/games">View upcoming games (JSON)</a> | <a href="/refresh">Force refresh</a></p>
  </body>
</html>
`, scheme, host)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func main() {
	godotenv.Load()

	cfg = config.Load()

	go refreshSchedule()

	c := cron.New()
	c.AddFunc(cfg.CacheRefreshInterval, func() {
		refreshSchedule()
	})
	c.Start()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calendar.ics", calendarHandler)
	http.HandleFunc("/games", gamesHandler)
	http.HandleFunc("/refresh", refreshHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Partizan ICS Calendar Server running on port %s", cfg.Port)
	log.Printf("Calendar URL: http://localhost:%s/calendar.ics", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal(err)
	}
}
