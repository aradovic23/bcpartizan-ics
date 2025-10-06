package ics

import (
	"fmt"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/milicacurcic/bcpartizan-ics/config"
	"github.com/milicacurcic/bcpartizan-ics/types"
)

func GenerateCalendar(games []types.Game, cfg *config.Config) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)

	for _, game := range games {
		event := cal.AddEvent(fmt.Sprintf("%s-%s-%s@partizan-ics", game.Date, game.HomeTeam, game.AwayTeam))

		dateTime := game.Date + " " + game.Time
		startTime, err := time.Parse("2006-01-02 15:04", dateTime)
		if err != nil {
			continue
		}

		endTime := startTime.Add(time.Duration(cfg.DefaultGameDuration) * time.Hour)

		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(startTime)
		event.SetEndAt(endTime)

		title := fmt.Sprintf("%s - %s vs %s", game.Competition, game.HomeTeam, game.AwayTeam)
		event.SetSummary(title)

		var description strings.Builder
		description.WriteString(game.Competition)
		description.WriteString("\n")
		description.WriteString(game.HomeTeam)
		description.WriteString(" vs ")
		description.WriteString(game.AwayTeam)
		if game.Round != "" {
			description.WriteString("\n")
			description.WriteString(game.Round)
		}
		if game.Venue != "" || game.Location != "" {
			description.WriteString("\nVenue: ")
			if game.Location != "" {
				description.WriteString(game.Location)
			} else {
				description.WriteString(game.Venue)
			}
		}
		event.SetDescription(description.String())

		if game.Location != "" {
			event.SetLocation(game.Location)
		} else if game.Venue != "" {
			event.SetLocation(game.Venue)
		}

		event.SetStatus(ics.ObjectStatusConfirmed)

		alarm30 := event.AddAlarm()
		alarm30.SetAction(ics.ActionDisplay)
		alarm30.SetTrigger("-PT30M")
		alarm30.SetDescription(fmt.Sprintf("%s starts in 30 minutes", title))

		alarm5 := event.AddAlarm()
		alarm5.SetAction(ics.ActionDisplay)
		alarm5.SetTrigger("-PT5M")
		alarm5.SetDescription(fmt.Sprintf("%s starts in 5 minutes", title))
	}

	return cal.Serialize(), nil
}
