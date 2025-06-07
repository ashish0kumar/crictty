package app

import (
	"crictty/internal/cricbuzz"
	"crictty/internal/models"
	"fmt"
	"time"
)

// App represents the main application structure
type App struct {
	client  *cricbuzz.Client
	Matches []models.MatchInfo
}

// New initializes a new App instance with all live matches
func New() (*App, error) {
	client := cricbuzz.NewClient()
	matches, err := client.GetAllLiveMatches()
	if err != nil {
		return nil, fmt.Errorf("failed to get live matches: %v", err)
	}

	return &App{
		client:  client,
		Matches: matches,
	}, nil
}

// NewWithMatchID initializes a new App instance with a specific match ID
func NewWithMatchID(matchID uint32) (*App, error) {
	client := cricbuzz.NewClient()
	matchInfo, err := client.GetMatchInfo(matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match info: %v", err)
	}

	shortName := fmt.Sprintf("%s vs %s",
		matchInfo.CricbuzzInfo.MatchHeader.Team1.ShortName,
		matchInfo.CricbuzzInfo.MatchHeader.Team2.ShortName)
	matchInfo.MatchShortName = shortName

	return &App{
		client:  client,
		Matches: []models.MatchInfo{matchInfo},
	}, nil
}

// UpdateMatches updates the matches in the App instance
func (a *App) UpdateMatches() error {
	if len(a.Matches) == 1 {
		// Single match mode -> update the specific match
		matchInfo, err := a.client.GetMatchInfo(a.Matches[0].CricbuzzMatchID)
		if err != nil {
			return err
		}
		matchInfo.MatchShortName = a.Matches[0].MatchShortName
		matchInfo.LastUpdated = time.Now()
		a.Matches[0] = matchInfo
	} else {
		// Multiple matches mode -> refresh all live matches
		matches, err := a.client.GetAllLiveMatches()
		if err != nil {
			return err
		}
		for i := range matches {
			matches[i].LastUpdated = time.Now()
		}
		a.Matches = matches
	}
	return nil
}

// GetMatchNames returns a slice of match names formatted for display
func (a *App) GetMatchNames() []string {
	names := make([]string, len(a.Matches))
	for i, match := range a.Matches {
		names[i] = fmt.Sprintf("%s - %s",
			match.MatchShortName,
			match.CricbuzzInfo.MatchHeader.MatchFormat)
	}
	return names
}
