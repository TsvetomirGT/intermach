package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const logoArt = `
  _       __            __  ___           __
 (_)___  / /____  _____/  |/  /___ ______/ /_
/ / __ \/ __/ _ \/ ___/ /|_/ / __ '/ ___/ __ \
/ / / / / /_/  __/ /  / /  / / /_/ / /__/ / / /
/_/_/ /_/\__/\___/_/  /_/  /_/\__,_/\___/_/ /_/
`

var funnyMessages = []string{
	"Waking up your router...",
	"Negotiating with the internet gods...",
	"Hunting for a fast server...",
	"Convincing packets to move faster...",
	"Bribing the DNS servers...",
}

func (m *Model) View() string {
	switch m.state {
	case stateInit:
		return m.viewInit()
	case stateFindingServer:
		return m.viewFinding()
	case stateDownload:
		return m.viewDownload()
	case stateUpload:
		return m.viewUpload()
	case stateFetchISP:
		return m.viewFetchISP()
	case stateDone:
		return m.viewDone()
	case stateError:
		return m.viewError()
	}
	return ""
}

func (m *Model) viewInit() string {
	logo := styleTitle.Render(logoArt)
	msg := styleBlue.Render("  Initializing internet speed tester...")
	quit := styleDim.Render("\n  press q to quit")
	return logo + "\n" + msg + quit
}

func (m *Model) viewFinding() string {
	logo := styleTitle.Render(logoArt)
	spin := m.spinner.View()
	msg := spin + " " + styleBlue.Render("Hunting for the nearest speed test server...")
	quit := styleDim.Render("\n  press q to quit")
	return logo + "\n  " + msg + quit
}

func (m *Model) viewDownload() string {
	logo := styleTitle.Render(logoArt)
	dl := m.liveDL

	bar := miniBar(dl, 1000, 20)
	speed := speedColor(dl).Render(fmt.Sprintf("%.1f Mbps", dl))
	line := fmt.Sprintf("  %s  ↓ Download   %s  %s",
		m.spinner.View(),
		speed,
		bar,
	)

	note := styleDim.Render("  ↑ Upload test follows...")
	quit := styleDim.Render("\n  press q to quit")
	return logo + "\n" + line + "\n" + note + quit
}

func (m *Model) viewUpload() string {
	logo := styleTitle.Render(logoArt)
	dl := m.liveDL
	ul := m.liveUL

	dlBar := miniBar(dl, 1000, 20)
	ulBar := miniBar(ul, 1000, 20)

	dlSpeed := speedColor(dl).Render(fmt.Sprintf("%.1f Mbps", dl))
	ulSpeed := speedColor(ul).Render(fmt.Sprintf("%.1f Mbps", ul))

	dlLine := fmt.Sprintf("  ✓ Download   %s  %s", dlSpeed, dlBar)
	ulLine := fmt.Sprintf("  %s ↑ Upload     %s  %s", m.spinner.View(), ulSpeed, ulBar)

	quit := styleDim.Render("\n  press q to quit")
	return logo + "\n" + dlLine + "\n" + ulLine + quit
}

func (m *Model) viewFetchISP() string {
	logo := styleTitle.Render(logoArt)
	spin := m.spinner.View()
	msg := spin + " " + styleBlue.Render("Asking the internet who you are...")
	quit := styleDim.Render("\n  press q to quit")
	return logo + "\n  " + msg + quit
}

func (m *Model) viewDone() string {
	r := m.result

	width := 44

	sep := styleLabel.Render(strings.Repeat("─", width))

	// Header
	header := styleTitle.Render("  intermach  speed results")

	// ISP info
	ispLine := ""
	ipLine := ""
	if m.isp != nil {
		ispLine = fmt.Sprintf("  %-12s %s", styleLabel.Render("ISP"), styleValue.Render(m.isp.Org))
		ipLine = fmt.Sprintf("  %-12s %s", styleLabel.Render("IP"), styleValue.Render(m.isp.IP))
	}
	serverLine := fmt.Sprintf("  %-12s %s",
		styleLabel.Render("Server"),
		styleValue.Render(fmt.Sprintf("%s, %s (%.0f ms)", r.ServerName, r.ServerCountry, r.LatencyMs)),
	)

	// Speeds
	maxSpeed := 1000.0
	if r.DownloadMbps > r.UploadMbps {
		maxSpeed = r.DownloadMbps * 1.2
	} else {
		maxSpeed = r.UploadMbps * 1.2
	}
	if maxSpeed < 10 {
		maxSpeed = 10
	}

	dlBar := miniBar(r.DownloadMbps, maxSpeed, 8)
	ulBar := miniBar(r.UploadMbps, maxSpeed, 8)
	latBar := miniBar(clamp(150-r.LatencyMs, 0, 150), 150, 8) // invert: lower latency = more bar

	dlLine := fmt.Sprintf("  %s %-10s  %s  %s",
		styleGreen.Render("↓"),
		styleLabel.Render("Download"),
		speedColor(r.DownloadMbps).Render(fmt.Sprintf("%6.1f Mbps", r.DownloadMbps)),
		dlBar,
	)
	ulLine := fmt.Sprintf("  %s %-10s  %s  %s",
		styleBlue.Render("↑"),
		styleLabel.Render("Upload"),
		speedColor(r.UploadMbps).Render(fmt.Sprintf("%6.1f Mbps", r.UploadMbps)),
		ulBar,
	)
	latLine := fmt.Sprintf("  %s %-10s  %s  %s",
		stylePurple.Render("◈"),
		styleLabel.Render("Latency"),
		latencyColor(r.LatencyMs).Render(fmt.Sprintf("%6.0f ms  ", r.LatencyMs)),
		latBar,
	)

	// Use-case ratings
	streaming4K := rating(r.DownloadMbps >= 25, r.UploadMbps >= 0, r.LatencyMs <= 9999)
	streamingHD := rating(r.DownloadMbps >= 5, r.UploadMbps >= 0, r.LatencyMs <= 9999)
	gaming := rating(r.DownloadMbps >= 3, r.UploadMbps >= 1, r.LatencyMs <= 100)
	videoCalls := rating(r.DownloadMbps >= 2.5, r.UploadMbps >= 2.5, r.LatencyMs <= 150)

	useLine1 := fmt.Sprintf("  Streaming HD  %s    Streaming 4K  %s", streamingHD, streaming4K)
	useLine2 := fmt.Sprintf("  Gaming        %s    Video Calls   %s", gaming, videoCalls)

	// Commentary
	commentary := funnyComment(r.DownloadMbps)
	commentLines := wordWrap(commentary, width-4)
	commentBlock := ""
	for _, line := range commentLines {
		commentBlock += "\n  " + styleDim.Render(`"`+line+`"`)
	}

	quit := "\n\n" + styleDim.Render("  press q to exit")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPurple).
		Padding(0, 1).
		Render(strings.Join([]string{
			"",
			header,
			sep,
			ispLine,
			ipLine,
			serverLine,
			sep,
			dlLine,
			ulLine,
			latLine,
			sep,
			useLine1,
			useLine2,
			sep,
			commentBlock,
			"",
		}, "\n"))

	return "\n" + box + quit
}

func (m *Model) viewError() string {
	return styleRed.Render("\n  Error: "+m.err.Error()) +
		"\n" + styleDim.Render("  press q to exit")
}

func rating(conditions ...bool) string {
	for _, c := range conditions {
		if !c {
			return styleRed.Render("✗")
		}
	}
	return styleGreen.Render("✓")
}

func funnyComment(dl float64) string {
	switch {
	case dl >= 200:
		return "Your internet is faster than your reaction time. Nice."
	case dl >= 50:
		return "Solid. You could stream 4K and still have bandwidth left for regrets."
	case dl >= 15:
		return "Good enough for calls, bad enough for excuses."
	case dl >= 5:
		return "Your internet runs on hope and ambition. Consider an upgrade."
	default:
		return "Have you tried turning your router off and moving to a city?"
	}
}

func wordWrap(s string, width int) []string {
	words := strings.Fields(s)
	var lines []string
	current := ""
	for _, w := range words {
		if len(current)+len(w)+1 > width {
			if current != "" {
				lines = append(lines, current)
			}
			current = w
		} else {
			if current == "" {
				current = w
			} else {
				current += " " + w
			}
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

var stylePurple = lipgloss.NewStyle().Foreground(colorPurple)
