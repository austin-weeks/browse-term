package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/austin-weeks/browse-term/internal/config"
	"github.com/austin-weeks/browse-term/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ch := make(chan string, 1)
	go checkLatestVersion(ch)
	js, err := checkJSEnabled()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	config := config.LoadConfig()

	a := tui.New(js, config)
	p := tea.NewProgram(a, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}

	select {
	case msg := <-ch:
		fmt.Print(msg)
	default:
	}
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

func checkLatestVersion(ch chan<- string) {
	sprintErr := func(msg string) string {
		return fmt.Sprintf("Could not check for latest version. Error: %s\n", msg)
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		ch <- sprintErr("could not read build info")
		return
	}
	curVer := info.Main.Version
	// Silence in development
	if curVer == "(devel)" {
		return
	}

	type entry struct {
		Tag string `json:"name"`
	}
	resp, err := http.Get("https://api.github.com/repos/austin-weeks/browse-term/tags")
	if err != nil {
		ch <- sprintErr(err.Error())
		return
	}
	defer resp.Body.Close() // nolint
	p, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- sprintErr(err.Error())
		return
	}
	var entries []entry
	err = json.Unmarshal(p, &entries)
	if err != nil {
		ch <- sprintErr(err.Error())
		return
	}
	if len(entries) == 0 {
		ch <- sprintErr("no tags found")
		return
	}
	newest := entries[0].Tag

	if curVer != newest {
		var s strings.Builder
		s.WriteString("\nA new version of browse-term is available!\n\n")
		s.WriteString(fmt.Sprintf("%s -> %s\n\n", curVer, newest))
		s.WriteString("To update, run:\n")
		s.WriteString("  go install github.com/austin-weeks/browse-term@latest\n")
		ch <- s.String()
	}
}
