package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashish0kumar/crictty/internal/app"
	"github.com/ashish0kumar/crictty/internal/models"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const mainWidth = 65 // Width of the main content area, adjust as needed

// keyMap defines the key bindings for the application
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Tab   key.Binding
	Quit  key.Binding
}

// Define key bindings for navigation and actions
var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous match"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next match"),
	),
	Tab: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "switch batting/bowling"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type tickMsg time.Time

// Model represents the state of the application
type Model struct {
	app            *app.App
	selectedMatch  int
	currentInnings int
	showBowling    bool
	tickRate       int
	width          int
	height         int
}

// NewModel creates a new Model instance with the given app and tick rate
func NewModel(app *app.App, tickRate int) Model {
	return Model{
		app:            app,
		selectedMatch:  0,
		currentInnings: 0,
		showBowling:    false,
		tickRate:       tickRate,
	}
}

// Init initializes the model, setting up the initial state and starting the tick command
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(m.tickRate),
	)
}

// tickCmd returns a command that ticks at the specified rate and updates matches
func tickCmd(tickRate int) tea.Cmd {
	return tea.Tick(time.Duration(tickRate)*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles incoming messages and updates the model state accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle window size changes
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	// Handle key messages for navigation and actions
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			if m.selectedMatch > 0 {
				m.selectedMatch--
				m.currentInnings = 0
				m.showBowling = false
			}
		case key.Matches(msg, keys.Right):
			if m.selectedMatch < len(m.app.Matches)-1 {
				m.selectedMatch++
				m.currentInnings = 0
				m.showBowling = false
			}
		case key.Matches(msg, keys.Up):
			if m.currentInnings > 0 {
				m.currentInnings--
			}
		case key.Matches(msg, keys.Down):
			if m.selectedMatch < len(m.app.Matches) {
				match := m.app.Matches[m.selectedMatch]
				if m.currentInnings < len(match.Scorecard)-1 {
					m.currentInnings++
				}
			}
		case key.Matches(msg, keys.Tab):
			m.showBowling = !m.showBowling
		}

	// Handle tick messages to update matches
	case tickMsg:
		cmds = append(cmds, tea.Cmd(func() tea.Msg {
			if err := m.app.UpdateMatches(); err != nil {
				return err
			}
			return nil
		}))
	}

	// Always schedule the next tick unless quitting
	cmds = append(cmds, tickCmd(m.tickRate))
	return m, tea.Batch(cmds...)
}

// View renders the current state of the model as a string
func (m Model) View() string {
	// If no matches are available show not found message
	if len(m.app.Matches) == 0 {
		return m.renderNotFoundMessage()
	}

	var content strings.Builder

	// Match tabs
	if len(m.app.Matches) > 1 {
		var tabs []string
		for i, name := range m.app.GetMatchNames() {
			style := tabStyle
			if i == m.selectedMatch {
				style = activeTabStyle
			}
			tabs = append(tabs, style.Render(name))
		}
		content.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
		content.WriteString("\n")
	}

	// Current match info
	if m.selectedMatch < len(m.app.Matches) {
		match := m.app.Matches[m.selectedMatch]
		content.WriteString(m.renderMatchInfo(match))
  	var match_id = fmt.Sprintf("Match id : %d", match.CricbuzzMatchID)
  	content.WriteString(helpStyle.Render(match_id))
  	content.WriteString("\n")
	}

	// Help
	content.WriteString("\n")
	content.WriteString(helpStyle.Render("q: quit • ←→: matches • ↑↓: innings • b: batting/bowling"))

	return m.centerHorizontally(content.String())
}

// centerHorizontally centers the content horizontally in the terminal
func (m Model) centerHorizontally(content string) string {
	return lipgloss.NewStyle().
		MarginTop(1).
		Width(m.width).
		Align(lipgloss.Center).
		Render(content)
}

// renderNotFoundMessage renders a message when no live matches are found
func (m Model) renderNotFoundMessage() string {
	notFoundMessage := "\nNo live matches found at the moment :(\n\n" +
		"This could be due to:\n\n" +
		"• Temporary issues with the Cricbuzz API\n" +
		"• No matches currently being played\n" +
		"• Your internet connection\n\n" +
		"Please try again in a few moments.\n\n" +
		"Use the --match-id flag with a valid match ID from Cricbuzz to view a specific match.\n\n"

	notFoundMessage += helpStyle.Render("Press 'q' to quit\n")
	return m.styleNotFoundMessage(notFoundMessage)
}

// styleNotFoundMessage styles the "no matches" message
func (m Model) styleNotFoundMessage(content string) string {
	return tabStyle.
		Width(mainWidth).
		MarginTop((m.height-20)/2).
		MarginLeft((m.width-mainWidth)/2).
		Padding(0, 4).
		Render(content)
}

// renderMatchInfo renders the match information including teams, scores, and current innings
func (m Model) renderMatchInfo(match models.MatchInfo) string {
	var content strings.Builder

	// Match header
	if len(m.app.Matches) <= 1 {
		header := fmt.Sprintf("%s vs %s - %s",
			match.CricbuzzInfo.MatchHeader.Team1.ShortName,
			match.CricbuzzInfo.MatchHeader.Team2.ShortName,
			match.CricbuzzInfo.MatchHeader.MatchFormat)
		headerStyled := activeTabStyle.Align(lipgloss.Center)
		content.WriteString(headerStyled.Render(header))
		content.WriteString("\n")
	}

	// Team scores
	content.WriteString("\n")
	content.WriteString(m.renderTeamScores(match.CricbuzzInfo.Miniscore.MatchScoreDetails))
	content.WriteString("\n\n")

	// Current innings info
	miniscore := match.CricbuzzInfo.Miniscore
	content.WriteString(m.renderCurrentInnings(miniscore))
	content.WriteString("\n")

	// Scorecard with batting/bowling tabs
	if len(match.Scorecard) > 0 && m.currentInnings < len(match.Scorecard) {
		content.WriteString(m.renderCurrentInningsScorecard(match.Scorecard[m.currentInnings], m.currentInnings))
	}

	return content.String()
}

// renderTeamScores renders the scores of both teams in a match
func (m Model) renderTeamScores(scoreDetails models.MatchScoreDetails) string {
	var content strings.Builder

	if len(scoreDetails.InningsScoreList) == 0 {
		return ""
	}

	leftWidth := mainWidth / 2
	rightWidth := mainWidth / 2

	var leftSide, rightSide strings.Builder

	// First innings (left side)
	if len(scoreDetails.InningsScoreList) > 0 {
		innings := scoreDetails.InningsScoreList[0]
		scoreText := m.formatInningsScore(innings)
		leftSide.WriteString(scoreStyle.Render(scoreText))

		// Third innings below first innings
		if len(scoreDetails.InningsScoreList) > 2 {
			innings3 := scoreDetails.InningsScoreList[2]
			scoreText3 := m.formatInningsScore(innings3)
			leftSide.WriteString("\n")
			leftSide.WriteString(scoreStyle.Render(scoreText3))
		}
	}

	// Second innings (right side)
	if len(scoreDetails.InningsScoreList) > 1 {
		innings := scoreDetails.InningsScoreList[1]
		scoreText := m.formatInningsScore(innings)
		rightSide.WriteString(scoreStyle.Render(scoreText))

		// Fourth innings below second innings
		if len(scoreDetails.InningsScoreList) > 3 {
			innings4 := scoreDetails.InningsScoreList[3]
			scoreText4 := m.formatInningsScore(innings4)
			rightSide.WriteString("\n")
			rightSide.WriteString(scoreStyle.Render(scoreText4))
		}
	}

	// Style containers
	leftContainer := lipgloss.NewStyle().
		Width(leftWidth).
		Align(lipgloss.Left).
		Render(leftSide.String())

	rightContainer := lipgloss.NewStyle().
		Width(rightWidth).
		Align(lipgloss.Right).
		Render(rightSide.String())

	// Join horizontally
	teamScoresRow := lipgloss.JoinHorizontal(lipgloss.Top, leftContainer, rightContainer)

	// Center the row
	centeredRow := lipgloss.NewStyle().
		Width(mainWidth).
		Align(lipgloss.Center).
		Render(teamScoresRow)

	content.WriteString(centeredRow)
	return content.String()
}

// formatInningsScore formats the innings score for display
func (m Model) formatInningsScore(innings models.InningsScore) string {
	if innings.IsDeclared {
		return fmt.Sprintf("%s %d/%d D (%.1f)",
			innings.BatTeamName,
			innings.Score,
			innings.Wickets,
			innings.Overs)
	} else if innings.Wickets == 10 {
		return fmt.Sprintf("%s %d (%.1f)",
			innings.BatTeamName,
			innings.Score,
			innings.Overs)
	} else {
		return fmt.Sprintf("%s %d/%d (%.1f)",
			innings.BatTeamName,
			innings.Score,
			innings.Wickets,
			innings.Overs)
	}
}

// renderCurrentInnings renders the current innings information including batsmen and bowler details
func (m Model) renderCurrentInnings(miniscore models.CricbuzzMiniscore) string {
	var content strings.Builder

	// Show live if there is no status and the match is in progress
	if miniscore.Status == "" && miniscore.MatchScoreDetails.State == "In Progress" {
		content.WriteString(statusStyle.Render("1st Innings"))
		content.WriteString("\n\n\n")
	}

	// Match status
	if miniscore.Status != "" {
		statusStyled := statusStyle.Width(mainWidth).Align(lipgloss.Center)
		content.WriteString(statusStyled.Render(miniscore.Status))
		content.WriteString("\n\n\n")
	}

	var leftSide, rightSide strings.Builder

	// Left side - Current batsmen
	striker := fmt.Sprintf("%s %d(%d)*",
		miniscore.BatsmanStriker.BatName,
		miniscore.BatsmanStriker.BatRuns,
		miniscore.BatsmanStriker.BatBalls)
	leftSide.WriteString(scoreStyle.Render(striker))
	leftSide.WriteString("\n")

	nonStriker := fmt.Sprintf("%s %d(%d)",
		miniscore.BatsmanNonStriker.BatName,
		miniscore.BatsmanNonStriker.BatRuns,
		miniscore.BatsmanNonStriker.BatBalls)
	leftSide.WriteString(scoreStyle.Render(nonStriker))

	// Right side - Current bowler
	bowlerName := miniscore.BowlerStriker.BowlName
	rightSide.WriteString(scoreStyle.Render(bowlerName))
	rightSide.WriteString("\n")

	bowlerFigures := fmt.Sprintf("%d-%d (%.1f)",
		miniscore.BowlerStriker.BowlWkts,
		miniscore.BowlerStriker.BowlRuns,
		miniscore.BowlerStriker.BowlOvs)
	rightSide.WriteString(scoreStyle.Render(bowlerFigures))

	// Layout current batsmen and bowler
	leftWidth := mainWidth * 2 / 3
	rightWidth := mainWidth / 3

	leftContainer := lipgloss.NewStyle().
		Width(leftWidth).
		Align(lipgloss.Left).
		Render(leftSide.String())

	rightContainer := lipgloss.NewStyle().
		Width(rightWidth).
		Align(lipgloss.Right).
		Render(rightSide.String())

	// Join horizontally
	currentInningsRow := lipgloss.JoinHorizontal(lipgloss.Top, leftContainer, rightContainer)
	content.WriteString(currentInningsRow)

	return content.String()
}

// renderCurrentInningsScorecard renders the scorecard for the current innings
func (m Model) renderCurrentInningsScorecard(innings models.MatchInningsInfo, inningsNumber int) string {
	var content strings.Builder

	match := m.app.Matches[m.selectedMatch]

	// Display innings indicator based on match format
	inningsIndicator := m.renderInningsIndicator(inningsNumber, len(match.Scorecard))
	scorecardTabs := m.renderScorecardTabs()

	headerRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(30).Align(lipgloss.Left).Render(inningsIndicator),
		lipgloss.NewStyle().Width(mainWidth-30).Align(lipgloss.Right).Render(scorecardTabs),
	)

	content.WriteString(fmt.Sprintf("\n%s\n", headerRow))

	// Display the Batting or Bowling card based on the toggle
	if m.showBowling {
		if len(innings.BowlerDetails) > 0 {
			content.WriteString(m.renderBowlingCard(innings.BowlerDetails))
		} else {
			content.WriteString(statusStyle.Render("No bowling data available for this innings"))
		}
	} else {
		if len(innings.BatsmanDetails) > 0 {
			content.WriteString(m.renderBattingCard(innings.BatsmanDetails))
		} else {
			content.WriteString(statusStyle.Render("No batting data available for this innings"))
		}
	}
	content.WriteString("\n")

	return content.String()
}

func (m Model) renderInningsIndicator(currentInnings, totalInnings int) string {
	// Create the tabs
	var inningsTabs []string
	for i := 0; i < totalInnings; i++ {
		tabLabel := ""
		switch i {
		case 0:
			tabLabel = fmt.Sprintf("%dst", i+1)
		case 1:
			tabLabel = fmt.Sprintf("%dnd", i+1)
		case 2:
			tabLabel = fmt.Sprintf("%drd", i+1)
		default:
			tabLabel = fmt.Sprintf("%dth", i+1)
		}

		if i == currentInnings {
			inningsTabs = append(inningsTabs, activeTabStyle.Render(tabLabel))
		} else {
			inningsTabs = append(inningsTabs, tabStyle.Render(tabLabel))
		}
	}

	// Join tabs together
	return lipgloss.JoinHorizontal(lipgloss.Center, inningsTabs...)
}

// renderScorecardTabs renders the tabs for switching between batting and bowling scorecards
func (m Model) renderScorecardTabs() string {
	var battingTab, bowlingTab string

	if m.showBowling {
		battingTab = tabStyle.Render("Bat")
		bowlingTab = activeTabStyle.Render("Bowl")
	} else {
		battingTab = activeTabStyle.Render("Bat")
		bowlingTab = tabStyle.Render("Bowl")
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Center, battingTab, bowlingTab)
	return lipgloss.NewStyle().MarginBottom(1).Render(tabs)
}

// renderBattingCard renders the batting scoreboard for the current innings
func (m Model) renderBattingCard(batsmen []models.BatsmanInfo) string {
	if len(batsmen) == 0 {
		return ""
	}

	var content strings.Builder

	otherColumnsWidth := 24 // 5 + 5 + 3 + 3 + 8 for R, B, 4s, 6s, S/R
	nameWidth := mainWidth - otherColumnsWidth - 7

	// Dynamic header formatting
	headerFormat := fmt.Sprintf("%%-%ds %%5s %%4s %%4s %%3s %%8s", nameWidth)
	headerRow := fmt.Sprintf(headerFormat, "Batsman", "R", "B", "4s", "6s", "S/R")

	content.WriteString(tableHeaderStyle.Render(headerRow))
	content.WriteString("\n")

	// Separator line
	separator := strings.Repeat("─", mainWidth)
	content.WriteString(helpStyle.Render(separator))
	content.WriteString("\n")

	// Data rows with dynamic formatting
	rowFormat := fmt.Sprintf("%%-%ds %%5s %%4s %%4s %%3s %%8s", nameWidth)

	// Data rows
	for _, bat := range batsmen {
		// Check if player is out
		isOut := bat.Status != "" &&
			bat.Status != "not out" &&
			bat.Status != "*" &&
			!strings.Contains(strings.ToLower(bat.Status), "not out")

		// Player stats row
		nameRow := fmt.Sprintf(rowFormat,
			truncateString(bat.Name, nameWidth),
			bat.Runs,
			bat.Balls,
			bat.Fours,
			bat.Sixes,
			bat.StrikeRate)

		content.WriteString(rowStyle.Render(nameRow))
		content.WriteString("\n")

		// Dismissal info below name
		if isOut {
			dismissalInfo := strings.TrimSpace(bat.Status)
			dismissalRow := dismissalInfo
			dismissalStyle := lipgloss.NewStyle().
				Width(mainWidth).
				Align(lipgloss.Left).
				PaddingLeft(1).
				Foreground(lipgloss.Color("8"))
			content.WriteString(dismissalStyle.Render(dismissalRow))
			content.WriteString("\n")
		} else {
			dismissalRow := "not out"
			dismissalStyle := lipgloss.NewStyle().
				Width(mainWidth).
				Align(lipgloss.Left).
				PaddingLeft(1).
				Foreground(lipgloss.Color("8"))
			content.WriteString(dismissalStyle.Render(dismissalRow))
			content.WriteString("\n")
		}
	}

	return content.String()
}

// renderBowlingCard renders the bowling scoreboard for the current innings
func (m Model) renderBowlingCard(bowlers []models.BowlerInfo) string {
	if len(bowlers) == 0 {
		return ""
	}

	var content strings.Builder

	// Calculate dynamic name column width
	otherColumnsWidth := 24 // 5 + 4 + 4 + 3 + 8 for O, M, R, W, Econ
	nameWidth := mainWidth - otherColumnsWidth - 7

	// Dynamic header formatting
	headerFormat := fmt.Sprintf("%%-%ds %%5s %%4s %%4s %%3s %%8s", nameWidth)
	headerRow := fmt.Sprintf(headerFormat, "Bowler", "O", "M", "R", "W", "Econ")

	content.WriteString(tableHeaderStyle.Render(headerRow))
	content.WriteString("\n")

	// Separator line
	separator := strings.Repeat("─", mainWidth)
	content.WriteString(helpStyle.Render(separator))
	content.WriteString("\n")

	// Data rows with dynamic formatting
	rowFormat := fmt.Sprintf("%%-%ds %%5s %%4s %%4s %%3s %%8s", nameWidth)

	// Data rows
	for _, bowl := range bowlers {
		// Bowler stats row
		nameRow := fmt.Sprintf(rowFormat,
			truncateString(bowl.Name, nameWidth),
			bowl.Overs,
			bowl.Maidens,
			bowl.Runs,
			bowl.Wickets,
			bowl.Economy)

		content.WriteString(rowStyle.Render(nameRow))
		content.WriteString("\n")
	}

	return content.String()
}

// truncateString truncates a string to a maximum length and appends "..." if truncated
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
