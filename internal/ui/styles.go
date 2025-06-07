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
)
