<div align="center">

```
  _       __            __  ___           __
 (_)___  / /____  _____/  |/  /___ ______/ /_
/ / __ \/ __/ _ \/ ___/ /|_/ / __ '/ ___/ __ \
/ / / / / /_/  __/ /  / /  / / /_/ / /__/ / / /
/_/_/ /_/\__/\___/_/  /_/  /_/\__,_/\___/_/ /_/
```

**A beautiful terminal speed test app — built with Go & Bubble Tea** 🚀

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/platform-macOS-lightgrey?style=flat-square&logo=apple)
![License](https://img.shields.io/badge/license-MIT-purple?style=flat-square)

</div>

---

## ✨ Features

- 📡 **Real speed test** — powered by speedtest.net, automatically finds the nearest server
- ⚡ **Live readout** — watch your download & upload speeds tick up in real time
- 🌍 **ISP detection** — shows your public IP, ISP name, and test server location
- 🎨 **Beautiful TUI** — animated, color-coded results with progress bars
- 🎮 **Use-case ratings** — instantly know if your connection handles 4K streaming, gaming, and video calls
- 😄 **Funny commentary** — because your connection deserves an honest opinion

---

## 📸 Results Panel

```
╭──────────────────────────────────────────────╮
│                                              │
│    intermach  speed results                  │
│  ────────────────────────────────────────    │
│  ISP        AS8866 Vivacom                   │
│  IP         84.x.x.x                         │
│  Server     Sofia, BG (12 ms)                │
│  ────────────────────────────────────────    │
│  ↓ Download     95.3 Mbps  ████████          │
│  ↑ Upload       48.1 Mbps  ████░░░░          │
│  ◈ Latency        12 ms    ████████          │
│  ────────────────────────────────────────    │
│  Streaming HD  ✓    Streaming 4K  ✓          │
│  Gaming        ✓    Video Calls   ✓          │
│  ────────────────────────────────────────    │
│                                              │
│  "Solid. You could stream 4K and still       │
│   have bandwidth left for regrets."          │
│                                              │
╰──────────────────────────────────────────────╯
```

Color coding: 🟢 great &nbsp; 🟡 okay &nbsp; 🔴 poor

---

## 🚀 Installation

### Option 1 — Homebrew (recommended)

```bash
brew tap TsvetomirGT/intermach
brew install intermach
```

### Option 2 — Install globally from source

```bash
git clone https://github.com/tsvetomirgt/intermach.git
cd intermach
make install       # may prompt for sudo password
```

This builds the binary and places it in `/usr/local/bin` so you can run `intermach` from anywhere.

### Option 3 — Build locally

```bash
make build
./intermach
```

### Option 4 — Run without building

```bash
make run
```

---

## 🎮 Usage

Just run:

```bash
intermach
```

Press **`q`** at any time to quit.

---

## 📊 Use-case Thresholds

| Use Case | Min Download | Min Upload | Max Latency |
|---|---|---|---|
| 📺 Streaming HD | 5 Mbps | — | — |
| 🎬 Streaming 4K | 25 Mbps | — | — |
| 🎮 Gaming | 3 Mbps | 1 Mbps | 100 ms |
| 📹 Video Calls (HD) | 2.5 Mbps | 2.5 Mbps | 150 ms |

---

## 🛠️ Development

### Commands

| Command | Description |
|---|---|
| `make build` | Compile to `./intermach` |
| `make run` | Run without building |
| `make install` | Build + install to `/usr/local/bin` |
| `make clean` | Remove local binary |
| `go vet ./...` | Lint the code |

### Project Structure

```
intermach/
├── main.go                    Entry point
└── internal/
    ├── network/
    │   └── info.go            ISP & IP via ipinfo.io
    ├── speedtest/
    │   └── runner.go          Speed test runner with live callbacks
    └── tui/
        ├── model.go           Bubbletea state machine
        ├── styles.go          Lipgloss styles & color palette
        └── view.go            All views & results panel
```

### Tech Stack

| Concern | Library |
|---|---|
| 🫧 TUI framework | [`charmbracelet/bubbletea`](https://github.com/charmbracelet/bubbletea) |
| 🎨 TUI styling | [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) |
| 🔄 Spinner | [`charmbracelet/bubbles`](https://github.com/charmbracelet/bubbles) |
| 📡 Speed testing | [`showwin/speedtest-go`](https://github.com/showwin/speedtest-go) |
| 🌍 ISP detection | [ipinfo.io](https://ipinfo.io) (no API key needed) |

---

## 📋 Requirements

- macOS (arm64 or amd64)
- Go 1.21+
- Internet connection 😄

---

<div align="center">

Made with ❤️ and Go

</div>
