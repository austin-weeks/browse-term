![Latest Version](https://img.shields.io/github/v/tag/austin-weeks/browse-term?label=latest)
![CICD](https://github.com/austin-weeks/browse-term/actions/workflows/CICD.yaml/badge.svg)
![Go version](https://img.shields.io/github/go-mod/go-version/austin-weeks/browse-term)
[![Go Report Card](https://goreportcard.com/badge/github.com/austin-weeks/browse-term)](https://goreportcard.com/report/github.com/austin-weeks/browse-term)
[![Go Reference](https://pkg.go.dev/badge/github.com/austin-weeks/browse-term.svg)](https://pkg.go.dev/github.com/austin-weeks/browse-term)

# <picture><img height="28" alt="Static Badge" src="https://img.shields.io/badge/%24-%230ea5e9?style=flat-square"></picture> BrowseTerm

A minimal web browser for your terminal.

![Browse-Term Preview](./.github/images/preview.gif)

## Installation

### Prerequisites

- A [Nerd Font](https://www.nerdfonts.com/) (recommended)
- A Chromium browser installed (optional, for JavaScript execution)

### Install with Go

```bash
go install github.com/austin-weeks/browse-term@latest
```

## Usage

Start BrowseTerm:

```bash
browse-term
```

See the _help_ line at the bottom of the viewport for available keyboard navigation.

### JavaScript Execution

By default, BrowseTerm uses Chromium to execute JavaScript on page load. This allows you to view content from sites that rely heavily on JavaScript for rendering the initial page.

You can disable JavaScript with the `--no-js` flag:

```bash
browse-term --no-js
```

## Project Roadmap

- [x] HTML as rendered markdown
- [x] Navigation keybinds
- [x] Tree for page links
- [x] JavaScript execution on page load
- [ ] Theme Configuration
- [ ] CSS Support
- [ ] Image Rendering
- [ ] Prebuilt Binaries in Releases
