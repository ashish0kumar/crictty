package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ashish0kumar/crictty/internal/app"
	"github.com/ashish0kumar/crictty/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	tickRate int
	matchID  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crictty",
	Short: "Live cricket scores in your terminal",
	Long:  "Minimal TUI for viewing cricket scoreboards right in your terminal",
	RunE:  runCrictty,
}

// SetVersion sets the version of the application
func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() error {
	return rootCmd.Execute()
}

// init initializes the root command and its flags
func init() {
	rootCmd.Flags().IntVarP(&tickRate, "tick-rate", "t", 40000, "Sets match details refresh rate in milliseconds")
	rootCmd.Flags().StringVarP(&matchID, "match-id", "m", "0", "ID of the match to follow live")
}

// runCrictty is the main function that runs the application
func runCrictty(cmd *cobra.Command, args []string) error {
	// Validate matchID input
	if matchID != "0" && !isValidMatchID(matchID) {
		return fmt.Errorf("invalid match ID format")
	}

	// Hide cursor during loading
	fmt.Print("\033[?25l")       // Hide cursor
	defer fmt.Print("\033[?25h") // Show cursor when function exits

	// Show simple loading message
	fmt.Print("\nFetching the scoreboard")

	// Simple loading animation
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for _, r := range `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` {
					fmt.Printf("\rFetching the scoreboard %c", r)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	// Initialize the application
	var cricketApp *app.App
	var err error

	if matchID == "0" {
		cricketApp, err = app.New()
	} else {
		id, _ := strconv.ParseUint(matchID, 10, 32)
		cricketApp, err = app.NewWithMatchID(uint32(id))
	}

	// Stop loading animation
	done <- true
	fmt.Print("\r                                    \r") // Clear loading line

	if err != nil {
		return fmt.Errorf("failed to load: %v", err)
	}

	// Start main UI
	model := ui.NewModel(cricketApp, tickRate)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %v", err)
	}

	return nil
}

// isValidMatchID checks if the provided match ID is a valid numeric string
func isValidMatchID(id string) bool {
	_, err := strconv.ParseUint(id, 10, 32)
	return err == nil
}
