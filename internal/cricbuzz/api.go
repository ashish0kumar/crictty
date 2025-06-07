package cricbuzz

import (
	"crictty/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// URL constants for Cricbuzz API endpoints
const (
	CricbuzzURL               = "https://www.cricbuzz.com"
	CricbuzzMatchAPI          = "https://www.cricbuzz.com/api/cricket-match/commentary/"
	CricbuzzMatchScorecardAPI = "https://www.cricbuzz.com/api/html/cricket-scorecard/"
)

// Client represents the Cricbuzz API client
type Client struct {
	httpClient *http.Client
}

// NewClient initializes a new Cricbuzz API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

const requestInterval = 1 * time.Second

var lastRequest time.Time

// makeRequest performs an HTTP GET request to the specified URL with rate limiting
func (c *Client) makeRequest(url string) (*http.Response, error) {
	time.Sleep(time.Until(lastRequest.Add(requestInterval)))
	lastRequest = time.Now()
	return c.httpClient.Get(url)
}

// cleanHTML removes unnecessary HTML tags and attributes from the given HTML content
func (c *Client) cleanHTML(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}

	cleanText := htmlContent

	re1 := regexp.MustCompile(`<span[^>]*class="[^"]*"[^>]*>`)
	cleanText = re1.ReplaceAllString(cleanText, "")

	cleanText = strings.ReplaceAll(cleanText, "</span>", "")

	re2 := regexp.MustCompile(`<a[^>]*>`)
	cleanText = re2.ReplaceAllString(cleanText, "")
	cleanText = strings.ReplaceAll(cleanText, "</a>", "")

	cleanText = strings.ReplaceAll(cleanText, "<strong>", "")
	cleanText = strings.ReplaceAll(cleanText, "</strong>", "")

	re3 := regexp.MustCompile(`<[^>]*>`)
	cleanText = re3.ReplaceAllString(cleanText, "")

	cleanText = strings.TrimSpace(cleanText)
	cleanText = regexp.MustCompile(`\s+`).ReplaceAllString(cleanText, " ")

	return cleanText
}

// GetAllLiveMatches fetches all live matches from Cricbuzz
func (c *Client) GetAllLiveMatches() ([]models.MatchInfo, error) {
	// Fetch the Cricbuzz homepage to get live matches
	resp, err := c.makeRequest(CricbuzzURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cricbuzz homepage: %v", err)
	}
	defer resp.Body.Close()

	// Parse the HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find all live matches in the navigation menu
	var matches []models.MatchInfo
	doc.Find("nav.cb-mat-mnu a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "" || text == "MATCHES" {
			return
		}

		// Check if the match is live
		parts := strings.Split(text, "-")
		if len(parts) >= 2 && strings.TrimSpace(parts[1]) == "Live" {
			href, exists := s.Attr("href")
			if !exists {
				return
			}

			// Extract match ID from the href
			pathParts := strings.Split(href, "/")
			if len(pathParts) < 3 {
				return
			}

			// Convert match ID to uint32
			matchID, err := strconv.ParseUint(pathParts[2], 10, 32)
			if err != nil {
				return
			}

			// Fetch match info using the match ID
			matchInfo, err := c.GetMatchInfo(uint32(matchID))
			if err != nil {
				return
			}

			// Set the match short name and append to the matches slice
			matchInfo.MatchShortName = strings.TrimSpace(parts[0])
			matches = append(matches, matchInfo)
		}
	})

	return matches, nil
}

// GetMatchInfo fetches detailed match information for a given match ID
func (c *Client) GetMatchInfo(matchID uint32) (models.MatchInfo, error) {
	// Construct the URL for the match API
	url := fmt.Sprintf("%s%d", CricbuzzMatchAPI, matchID)
	resp, err := c.makeRequest(url)
	if err != nil {
		return models.MatchInfo{}, fmt.Errorf("failed to fetch match info: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	var cricbuzzJSON models.CricbuzzJSON
	if err := json.NewDecoder(resp.Body).Decode(&cricbuzzJSON); err != nil {
		return models.MatchInfo{}, fmt.Errorf("failed to decode JSON: %v", err)
	}

	// Check if the match is complete
	scorecard, err := c.GetScorecard(matchID)
	if err != nil {
		scorecard = []models.MatchInningsInfo{}
	}

	return models.MatchInfo{
		CricbuzzMatchID:      matchID,
		CricbuzzMatchAPILink: url,
		CricbuzzInfo:         cricbuzzJSON,
		Scorecard:            scorecard,
	}, nil
}

// GetScorecard fetches the scorecard for a given match ID
func (c *Client) GetScorecard(matchID uint32) ([]models.MatchInningsInfo, error) {
	// Construct the URL for the scorecard API
	url := fmt.Sprintf("%s%d", CricbuzzMatchScorecardAPI, matchID)
	resp, err := c.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scorecard: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scorecard HTML: %v", err)
	}

	var scorecard []models.MatchInningsInfo

	// Iterate through the innings sections (1 to 4)
	for i := 1; i <= 4; i++ {
		selector := fmt.Sprintf("div[id=\"innings_%d\"]", i)
		inningsDiv := doc.Find(selector)
		if inningsDiv.Length() == 0 {
			continue
		}

		innings := c.parseInningsInfo(inningsDiv)
		scorecard = append(scorecard, innings)
	}

	return scorecard, nil
}

// parseInningsInfo parses the innings information from the given goquery selection
func (c *Client) parseInningsInfo(inningsDiv *goquery.Selection) models.MatchInningsInfo {
	var innings models.MatchInningsInfo

	// Extract innings ID from the div ID
	inningsDiv.Find("div.cb-scrd-itms").Each(func(i int, s *goquery.Selection) {
		divs := s.Find("div")
		divCount := divs.Length()

		if divCount >= 6 {
			var firstCol, secondCol string
			if divCount > 0 {
				firstCol = strings.TrimSpace(divs.Eq(0).Text())
			}
			if divCount > 1 {
				secondCol = strings.TrimSpace(divs.Eq(1).Text())
			}

			// Check if this is batting or bowling data
			if c.isBattingRow(firstCol, secondCol, divCount) {

				// Parse batting data
				batsman := models.BatsmanInfo{}
				divs.Each(func(j int, div *goquery.Selection) {
					text := strings.TrimSpace(div.Text())
					html, _ := div.Html()

					switch j {
					case 0:
						batsman.Name = text
					case 1:
						batsman.Status = c.cleanHTML(html)
					case 2:
						batsman.Runs = c.cleanHTML(html)
					case 3:
						batsman.Balls = c.cleanHTML(html)
					case 4:
						batsman.Fours = c.cleanHTML(html)
					case 5:
						batsman.Sixes = c.cleanHTML(html)
					case 6:
						if j < divCount {
							batsman.StrikeRate = c.cleanHTML(html)
						}
					}
				})

				if batsman.Name != "" &&
					!strings.Contains(strings.ToLower(batsman.Name), "extras") &&
					!strings.Contains(strings.ToLower(batsman.Name), "total") &&
					!strings.Contains(strings.ToLower(batsman.Name), "fall of wickets") {
					innings.BatsmanDetails = append(innings.BatsmanDetails, batsman)
				}
			} else {

				// Parse bowling data
				bowler := models.BowlerInfo{}
				divs.Each(func(j int, div *goquery.Selection) {
					text := strings.TrimSpace(div.Text())
					html, _ := div.Html()

					switch j {
					case 0:
						bowler.Name = text
					case 1:
						bowler.Overs = c.cleanHTML(html)
					case 2:
						bowler.Maidens = c.cleanHTML(html)
					case 3:
						bowler.Runs = c.cleanHTML(html)
					case 4:
						bowler.Wickets = c.cleanHTML(html)
					case 5:
						bowler.NoBalls = c.cleanHTML(html)
					case 6:
						bowler.Wides = c.cleanHTML(html)
					case 7:
						if j < divCount {
							bowler.Economy = c.cleanHTML(html)
						}
					}
				})

				if bowler.Name != "" {
					innings.BowlerDetails = append(innings.BowlerDetails, bowler)
				}
			}
		}
	})

	return innings
}

// isBattingRow determines if the given row represents batting data based on its content
func (c *Client) isBattingRow(firstCol, secondCol string, divCount int) bool {
	// Common batting status indicators
	battingStatuses := []string{
		"not out", "c ", "b ", "lbw", "run out", "st ", "hit wicket",
		"obstructing", "handled", "timed out", "*",
	}

	// Common bowling indicators (overs format like "4.0", "10.2")
	oversPattern := regexp.MustCompile(`^\d+\.\d+$`)

	// If first column is empty, assume it's not batting
	if oversPattern.MatchString(secondCol) {
		return false
	}

	// Check if the second column contains any batting status
	for _, status := range battingStatuses {
		if strings.Contains(strings.ToLower(secondCol), status) {
			return true
		}
	}

	// Skip rows that are not relevant to batting
	skipRows := []string{
		"extras", "total", "fall of wickets", "bowler", "overs", "maidens",
		"runs", "wickets", "economy", "nb", "wd",
	}

	// Check if the first or second column contains any skip indicators
	for _, skip := range skipRows {
		if strings.Contains(strings.ToLower(firstCol), skip) ||
			strings.Contains(strings.ToLower(secondCol), skip) {
			return false
		}
	}

	// If there are 6 to 8 columns and the first column is not empty, assume it's batting
	if divCount >= 6 && divCount <= 8 && firstCol != "" {
		return true
	}

	return false
}
