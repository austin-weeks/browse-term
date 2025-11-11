package main

import (
	"log"

	"github.com/austin-weeks/browse-term/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO - tui should be restructured to use new structs rather than mutating things
// TODO - tui should be restructured to use non-interface types for sub-components - didn't need to be this abstract

func main() {
	a := tui.New()
	p := tea.NewProgram(a, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// Clear the screen

	// Check for latest version of app
	checkLatestVersion()
}

func checkLatestVersion() {
}
