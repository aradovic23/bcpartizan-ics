package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/milicacurcic/bcpartizan-ics/config"
	"github.com/milicacurcic/bcpartizan-ics/types"
)

var venueMap = map[string]string{
	"Partizan":         "Stark Arena, Belgrade, Serbia",
	"Crvena zvezda":    "Aleksandar Nikolic Hall, Belgrade, Serbia",
	"Real Madrid":      "WiZink Center, Madrid, Spain",
	"Barcelona":        "Palau Blaugrana, Barcelona, Spain",
	"Panathinaikos":    "OAKA, Athens, Greece",
	"Olympiacos":       "Peace and Friendship Stadium, Piraeus, Greece",
	"Fenerbahce":       "Ulker Sports Arena, Istanbul, Turkey",
	"Anadolu Efes":     "Sinan Erdem Dome, Istanbul, Turkey",
	"Maccabi Tel Aviv": "Menora Mivtachim Arena, Tel Aviv, Israel",
	"Zalgiris":         "Zalgirio Arena, Kaunas, Lithuania",
	"Bayern Munich":    "Audi Dome, Munich, Germany",
	"ALBA Berlin":      "Mercedes-Benz Arena, Berlin, Germany",
	"Milano":           "Mediolanum Forum, Milan, Italy",
	"Virtus Bologna":   "Segafredo Arena, Bologna, Italy",
	"Monaco":           "Salle Gaston Médecin, Monaco",
	"Paris":            "Adidas Arena, Paris, France",
	"Baskonia":         "Buesa Arena, Vitoria-Gasteiz, Spain",
	"Valencia":         "La Fonteta, Valencia, Spain",
}

func getVenueForTeam(teamName string) string {
	for team, venue := range venueMap {
		if strings.Contains(teamName, team) {
			return venue
		}
	}
	return teamName + " Arena"
}

func parseFlashScoreData(html, team string) []types.Game {
	var games []types.Game

	dataRegex := regexp.MustCompile(`data:\s*` + "`" + `([^` + "`" + `]+)` + "`")
	dataMatch := dataRegex.FindStringSubmatch(html)
	if len(dataMatch) < 2 {
		return games
	}

	dataString := dataMatch[1]
	gameBlocks := strings.Split(dataString, "~AA÷")

	for i := range gameBlocks {
		block := gameBlocks[i]
		if !strings.Contains(block, "PTZ") && !strings.Contains(block, team) {
			continue
		}

		fields := make(map[string]string)
		pairs := strings.SplitSeq(block, "¬")

		for pair := range pairs {
			parts := strings.Split(pair, "÷")
			if len(parts) == 2 {
				fields[parts[0]] = parts[1]
			}
		}

		var venue string
		if i < len(gameBlocks)-1 {
			nextBlock := gameBlocks[i+1]
			venueRegex := regexp.MustCompile(`AM÷([^¬]+)`)
			venueMatch := venueRegex.FindStringSubmatch(nextBlock)
			if len(venueMatch) > 1 {
				venue = strings.TrimSuffix(strings.Replace(venueMatch[1], "Neutral location - ", "", 1), ".")
			}
		}

		timestamp := fields["AD"]
		homeTeam := fields["AE"]
		awayTeam := fields["AF"]
		round := fields["ER"]
		competition := fields["ZA"]

		if (strings.Contains(homeTeam, team) || strings.Contains(awayTeam, team)) && timestamp != "" {
			ts, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				continue
			}

			date := time.Unix(ts, 0)

			if venue == "" || strings.Contains(venue, "TBD") {
				venue = getVenueForTeam(homeTeam)
			}

			if competition == "" {
				competition = "Euroleague"
			}
			if homeTeam == "" {
				homeTeam = "Unknown"
			}
			if awayTeam == "" {
				awayTeam = "Unknown"
			}

			games = append(games, types.Game{
				Competition: competition,
				HomeTeam:    homeTeam,
				AwayTeam:    awayTeam,
				Date:        date.Format("2006-01-02"),
				Time:        date.Format("15:04"),
				Venue:       venue,
				Location:    venue,
				Round:       round,
				Source:      "flashscore",
			})
		}
	}

	return games
}

func fetchFlashScoreSchedule(url, competition, team string) []types.Game {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v", competition, err)
		return []types.Game{}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching %s schedule from FlashScore: %v", competition, err)
		return []types.Game{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for %s: %v", competition, err)
		return []types.Game{}
	}

	games := parseFlashScoreData(string(body), team)
	for i := range games {
		games[i].Competition = competition
		games[i].Source = "flashscore"
	}

	return games
}

type EuroleagueResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Date string `json:"date"`
		Home struct {
			Name            string `json:"name"`
			AbbreviatedName string `json:"abbreviatedName"`
		} `json:"home"`
		Away struct {
			Name            string `json:"name"`
			AbbreviatedName string `json:"abbreviatedName"`
		} `json:"away"`
		Venue struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"venue"`
		Round struct {
			Name string `json:"name"`
		} `json:"round"`
	} `json:"data"`
}

func fetchEuroleagueSchedule(cfg *config.Config) []types.Game {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("%s?teamCode=%s", cfg.EuroleagueAPIURL, cfg.TeamCodeEuroleague)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating Euroleague request: %v", err)
		return []types.Game{}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching Euroleague API: %v", err)
		return []types.Game{}
	}
	defer resp.Body.Close()

	var apiResp EuroleagueResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Printf("Error decoding Euroleague response: %v", err)
		return []types.Game{}
	}

	if apiResp.Status != "success" || apiResp.Data == nil {
		log.Printf("Euroleague API returned no data")
		return []types.Game{}
	}

	now := time.Now()
	var games []types.Game

	for _, game := range apiResp.Data {
		gameDate, err := time.Parse(time.RFC3339, game.Date)
		if err != nil {
			continue
		}

		if gameDate.Before(now) {
			continue
		}

		homeTeam := game.Home.AbbreviatedName
		if homeTeam == "" {
			homeTeam = game.Home.Name
		}

		awayTeam := game.Away.AbbreviatedName
		if awayTeam == "" {
			awayTeam = game.Away.Name
		}

		venue := ""
		if game.Venue.Name != "" && game.Venue.Address != "" {
			venue = fmt.Sprintf("%s, %s", game.Venue.Name, game.Venue.Address)
		} else if game.Venue.Name != "" {
			venue = game.Venue.Name
		} else {
			venue = getVenueForTeam(game.Home.Name)
		}

		games = append(games, types.Game{
			Competition: "Euroleague",
			HomeTeam:    homeTeam,
			AwayTeam:    awayTeam,
			Date:        gameDate.Format("2006-01-02"),
			Time:        gameDate.Format("15:04"),
			Venue:       venue,
			Location:    venue,
			Round:       game.Round.Name,
			Source:      "euroleague-api",
		})
	}

	return games
}

func deduplicateGames(games []types.Game) []types.Game {
	seen := make(map[string]types.Game)

	for _, game := range games {
		key := fmt.Sprintf("%s-%s-%s-%s", game.Date, game.Time, game.HomeTeam, game.AwayTeam)
		if _, exists := seen[key]; !exists {
			seen[key] = game
		}
	}

	var result []types.Game
	for _, game := range seen {
		result = append(result, game)
	}

	return result
}

func fetchABALeagueSchedule(cfg *config.Config) []types.Game {
	fixtures := fetchFlashScoreSchedule(cfg.FlashscoreABAFixturesURL, "ABA League", cfg.Team)
	results := fetchFlashScoreSchedule(cfg.FlashscoreABAResultsURL, "ABA League", cfg.Team)

	allGames := append(fixtures, results...)
	now := time.Now()

	var futureGames []types.Game
	for _, game := range allGames {
		gameDateTime := game.Date + " " + game.Time
		gameDate, err := time.Parse("2006-01-02 15:04", gameDateTime)
		if err != nil {
			continue
		}

		if gameDate.After(now) || gameDate.Equal(now) {
			futureGames = append(futureGames, game)
		}
	}

	return deduplicateGames(futureGames)
}

func FetchAllSchedules(cfg *config.Config) []types.Game {
	euroleagueGames := fetchEuroleagueSchedule(cfg)
	abaGames := fetchABALeagueSchedule(cfg)

	allGames := append(euroleagueGames, abaGames...)

	sort.Slice(allGames, func(i, j int) bool {
		dateTimeA := allGames[i].Date + " " + allGames[i].Time
		dateTimeB := allGames[j].Date + " " + allGames[j].Time

		dateA, errA := time.Parse("2006-01-02 15:04", dateTimeA)
		dateB, errB := time.Parse("2006-01-02 15:04", dateTimeB)

		if errA != nil || errB != nil {
			return false
		}

		return dateA.Before(dateB)
	})

	return allGames
}

func GetMockSchedule() []types.Game {
	now := time.Now()
	var games []types.Game

	opponents := []string{"Real Madrid", "Barcelona", "Olimpia Milano", "Fenerbahce", "Crvena Zvezda", "Maccabi", "Bayern Munich", "Zalgiris"}

	for i := 0; i < 10; i++ {
		gameDate := now.AddDate(0, 0, i*7)
		opponent := opponents[i%len(opponents)]
		isHome := i%2 == 0

		homeTeam := "Partizan"
		awayTeam := opponent
		venue := "Stark Arena, Belgrade, Serbia"

		if !isHome {
			homeTeam = opponent
			awayTeam = "Partizan"
			venue = opponent + " Arena"
		}

		competition := "ABA League"
		if i%3 == 0 {
			competition = "Euroleague"
		}

		games = append(games, types.Game{
			Competition: competition,
			HomeTeam:    homeTeam,
			AwayTeam:    awayTeam,
			Date:        gameDate.Format("2006-01-02"),
			Time:        "20:00",
			Venue:       venue,
			Location:    venue,
			Source:      "mock",
		})
	}

	return games
}
