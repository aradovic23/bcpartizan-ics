package types

type Game struct {
	Competition string `json:"competition"`
	HomeTeam    string `json:"homeTeam"`
	AwayTeam    string `json:"awayTeam"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Venue       string `json:"venue"`
	Location    string `json:"location"`
	Round       string `json:"round,omitempty"`
	Source      string `json:"source"`
}
