package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles for the TUI components
var (
	tabStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 1)

	activeTabStyle = tabStyle.
			BorderForeground(lipgloss.Color("15")).
			Bold(true)

	scoreStyle = lipgloss.NewStyle().
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Italic(true)

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Padding(0, 1)

	rowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Padding(0, 1).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	// Progress and loading styles
	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	connectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	// Team and status color styles
	team1Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")).
			Bold(true)

	team2Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true)

	liveMatchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)

	completedMatchStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))

	boundaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Bold(true)

	wicketStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
)
