package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorGreen  = lipgloss.Color("#00D26A")
	colorYellow = lipgloss.Color("#FFD700")
	colorRed    = lipgloss.Color("#FF4444")
	colorBlue   = lipgloss.Color("#4FC3F7")
	colorPurple = lipgloss.Color("#CE93D8")
	colorGray   = lipgloss.Color("#666666")
	colorWhite  = lipgloss.Color("#FFFFFF")

	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple)

	styleLabel = lipgloss.NewStyle().
			Foreground(colorGray)

	styleValue = lipgloss.NewStyle().
			Foreground(colorWhite).
			Bold(true)

	styleGreen = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	styleYellow = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	styleRed = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	styleBlue = lipgloss.NewStyle().
			Foreground(colorBlue)

	styleDim = lipgloss.NewStyle().
			Foreground(colorGray)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPurple).
			Padding(0, 2)

	styleStatusBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorBlue).
			Padding(0, 1)
)

func speedColor(mbps float64) lipgloss.Style {
	switch {
	case mbps >= 50:
		return styleGreen
	case mbps >= 10:
		return styleYellow
	default:
		return styleRed
	}
}

func latencyColor(ms float64) lipgloss.Style {
	switch {
	case ms <= 30:
		return styleGreen
	case ms <= 100:
		return styleYellow
	default:
		return styleRed
	}
}

func miniBar(val, max float64, width int) string {
	if max == 0 {
		return ""
	}
	filled := int(val / max * float64(width))
	if filled > width {
		filled = width
	}
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += styleGreen.Render("█")
		} else {
			bar += styleDim.Render("░")
		}
	}
	return bar
}
