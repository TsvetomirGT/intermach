# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

**intermach** is a terminal TUI speed test app written in Go. It measures your download/upload speed, latency, and ISP info, displaying results in an animated terminal interface.

## Commands

- Build: `make build` — produces `./intermach` binary
- Run: `make run` — `go run .`
- Install: `make install` — builds and copies binary to `/usr/local/bin/intermach` (may require `sudo`)
- Clean: `make clean` — removes the local `./intermach` binary
- Vet: `go vet ./...`

## Architecture

```
intermach/
├── main.go                    Entry point — creates and runs the bubbletea program
└── internal/
    ├── network/
    │   └── info.go            Fetches public IP and ISP via GET https://ipinfo.io/json
    ├── speedtest/
    │   └── runner.go          Wraps showwin/speedtest-go; returns Result struct with live progress callbacks
    └── tui/
        ├── model.go           Bubbletea Model, Msg types, state machine, tea.Cmd wiring
        ├── styles.go          Lipgloss color palette, border styles, miniBar helper
        └── view.go            View() rendering — all ASCII art, animated states, results panel
```

### TUI State Machine

```
stateInit → stateFindingServer → stateDownload → stateUpload → stateFetchISP → stateDone
```

### Stack

| Concern | Library |
|---|---|
| TUI framework | `github.com/charmbracelet/bubbletea` |
| TUI styling | `github.com/charmbracelet/lipgloss` |
| TUI components | `github.com/charmbracelet/bubbles` (spinner) |
| Speed testing | `github.com/showwin/speedtest-go` |
| ISP detection | `https://ipinfo.io/json` (no API key, 50k/month free) |

### Live Progress

The speed test runs inside a `tea.Cmd` goroutine. Progress (`float64` Mbps values) is streamed to the model via a buffered channel. A 100ms ticker drains the channel and triggers re-renders, giving a live speed readout during download and upload phases.
