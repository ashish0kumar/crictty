package models

import "time"

// BowlerInfo contains bowling statistics for a player in a single innings
type BowlerInfo struct {
	Name    string
	Overs   string
	Maidens string
	Runs    string
	Wickets string
	NoBalls string
	Wides   string
	Economy string
}

// BatsmanInfo contains batting statistics for a player in a single innings
type BatsmanInfo struct {
	Name       string
	Status     string
	Runs       string
	Balls      string
	Fours      string
	Sixes      string
	StrikeRate string
}

// MatchInningsInfo holds all batting and bowling details for an innings
type MatchInningsInfo struct {
	BatsmanDetails []BatsmanInfo
	YetToBat       string
	BowlerDetails  []BowlerInfo
}

// CricbuzzMiniscore contains live match summary information.
type CricbuzzMiniscore struct {
	InningsID         uint32            `json:"inningsId"`
	BatsmanStriker    Batsman           `json:"batsmanStriker"`
	BatsmanNonStriker Batsman           `json:"batsmanNonStriker"`
	BowlerStriker     Bowler            `json:"bowlerStriker"`
	BowlerNonStriker  Bowler            `json:"bowlerNonStriker"`
	Overs             float32           `json:"overs"`
	RecentOvsStats    string            `json:"recentOvsStats"`
	CurrentRunRate    float32           `json:"currentRunRate"`
	RequiredRunRate   float32           `json:"requiredRunRate"`
	LastWicket        *string           `json:"lastWicket"`
	MatchScoreDetails MatchScoreDetails `json:"matchScoreDetails"`
	OversRem          *float32          `json:"oversRem"`
	Status            string            `json:"status"`
}

// Batsman represents batting stats for a player as received from the API
type Batsman struct {
	BatBalls      uint32  `json:"batBalls"`
	BatDots       uint32  `json:"batDots"`
	BatFours      uint32  `json:"batFours"`
	BatID         uint32  `json:"batId"`
	BatName       string  `json:"batName"`
	BatMins       uint32  `json:"batMins"`
	BatSixes      uint32  `json:"batSixes"`
	BatStrikeRate float32 `json:"batStrikeRate"`
	BatRuns       uint32  `json:"batRuns"`
}

// Bowler represents bowling stats for a player as received from the API
type Bowler struct {
	BowlID      uint32  `json:"bowlId"`
	BowlName    string  `json:"bowlName"`
	BowlMaidens uint32  `json:"bowlMaidens"`
	BowlNoballs uint32  `json:"bowlNoballs"`
	BowlOvs     float32 `json:"bowlOvs"`
	BowlRuns    uint32  `json:"bowlRuns"`
	BowlWides   uint32  `json:"bowlWides"`
	BowlWkts    uint32  `json:"bowlWkts"`
	BowlEcon    float32 `json:"bowlEcon"`
}

// MatchScoreDetails contains overall match score and innings summary
type MatchScoreDetails struct {
	MatchID          uint32         `json:"matchId"`
	MatchFormat      string         `json:"matchFormat"`
	State            string         `json:"state"`
	CustomStatus     string         `json:"customStatus"`
	InningsScoreList []InningsScore `json:"inningsScoreList"`
}

// InningsScore contains summary of runs, wickets, and overs for an innings
type InningsScore struct {
	InningsID   uint32  `json:"inningsId"`
	BatTeamID   uint32  `json:"batTeamId"`
	BatTeamName string  `json:"batTeamName"`
	Score       uint32  `json:"score"`
	Wickets     uint32  `json:"wickets"`
	Overs       float32 `json:"overs"`
	IsDeclared  bool    `json:"isDeclared"`
	IsFollowOn  bool    `json:"isFollowOn"`
}

// MatchHeader contains metadata about the match
type MatchHeader struct {
	MatchID                uint32  `json:"matchId"`
	MatchDescription       string  `json:"matchDescription"`
	MatchFormat            string  `json:"matchFormat"`
	MatchType              string  `json:"matchType"`
	Complete               bool    `json:"complete"`
	Domestic               bool    `json:"domestic"`
	MatchStartTimestamp    uint64  `json:"matchStartTimestamp"`
	MatchCompleteTimestamp uint64  `json:"matchCompleteTimestamp"`
	DayNight               *bool   `json:"dayNight"`
	Year                   uint32  `json:"year"`
	DayNumber              *uint32 `json:"dayNumber"`
	State                  string  `json:"state"`
	Status                 string  `json:"status"`
	Team1                  Team    `json:"team1"`
	Team2                  Team    `json:"team2"`
	SeriesDesc             string  `json:"seriesDesc"`
	SeriesID               uint32  `json:"seriesId"`
	SeriesName             string  `json:"seriesName"`
}

// Team contains team identification and names
type Team struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

// CricbuzzJSON contains match header, miniscore, and page info
type CricbuzzJSON struct {
	MatchHeader MatchHeader       `json:"matchHeader"`
	Miniscore   CricbuzzMiniscore `json:"miniscore"`
	Page        string            `json:"page"`
}

// MatchInfo contains match metadata, live data, and scorecard
type MatchInfo struct {
	MatchShortName       string
	CricbuzzMatchID      uint32
	CricbuzzMatchAPILink string
	CricbuzzInfo         CricbuzzJSON
	Scorecard            []MatchInningsInfo
	LastUpdated          time.Time
}
