package main

import (
	"fmt"
	"log"
	"os"

	"github.com/austin-weeks/browse-term/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO - tui should be restructured to use new structs rather than mutating things
// TODO - tui should be restructured to use non-interface types for sub-components - didn't need to be this abstract

func main() {
	js, err := checkJSEnabled()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	a := tui.New(js)
	p := tea.NewProgram(a, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// Check for latest version of app
	checkLatestVersion()
}

func checkJSEnabled() (bool, error) {
	enabled := true
	for _, flag := range os.Args[1:] {
		switch flag {
		case "--no-js":
			enabled = false
		default:
			return false, fmt.Errorf("unknown option %s", flag)
		}
	}
	return enabled, nil
}

func checkLatestVersion() {
}
