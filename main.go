package main

import (
	"log"

	"github.com/austin-weeks/browse-term/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	a := tui.New()
	p := tea.NewProgram(a)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
