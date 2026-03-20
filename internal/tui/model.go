package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tsvetomirgt/intermach/internal/network"
	"github.com/tsvetomirgt/intermach/internal/speedtest"
)

type state int

const (
	stateInit state = iota
	stateFindingServer
	stateDownload
	stateUpload
	stateFetchISP
	stateDone
	stateError
)

// Messages
type msgStartTest struct{}
type msgTick struct{}
type msgDLProgress struct{ mbps float64 }
type msgULProgress struct{ mbps float64 }
type msgResult struct{ result *speedtest.Result }
type msgISPDone struct{ info *network.IPInfo }
type msgError struct{ err error }

// progressUpdate is sent through the tea.Cmd channel mechanism.
type progressUpdate struct {
	dl *float64 // non-nil when download progress
	ul *float64 // non-nil when upload progress
}

type Model struct {
	state   state
	spinner spinner.Model
	err     error

	liveDL float64
	liveUL float64

	result *speedtest.Result
	isp    *network.IPInfo

	// channel for streaming progress from the speed test goroutine
	progressCh chan progressUpdate
}

func New() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styleBlue
	return &Model{
		state:      stateInit,
		spinner:    s,
		progressCh: make(chan progressUpdate, 64),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Tick(400*time.Millisecond, func(_ time.Time) tea.Msg {
			return msgStartTest{}
		}),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(_ time.Time) tea.Msg {
		return msgTick{}
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case msgTick:
		// Drain progress channel and pick the latest value
		drained := false
		for {
			select {
			case p := <-m.progressCh:
				if p.dl != nil {
					m.liveDL = *p.dl
					m.state = stateDownload
				}
				if p.ul != nil {
					m.liveUL = *p.ul
					m.state = stateUpload
				}
				drained = true
			default:
				goto done
			}
		}
	done:
		_ = drained
		return m, tickCmd()

	case msgStartTest:
		m.state = stateFindingServer
		return m, tea.Batch(m.launchSpeedTest(), tickCmd())

	case msgResult:
		m.result = msg.result
		m.state = stateFetchISP
		return m, m.fetchISP()

	case msgISPDone:
		m.isp = msg.info
		m.state = stateDone
		return m, nil

	case msgError:
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	return m, nil
}

func (m *Model) launchSpeedTest() tea.Cmd {
	ch := m.progressCh
	return func() tea.Msg {
		result, err := speedtest.Run(
			func(mbps float64) {
				v := mbps
				select {
				case ch <- progressUpdate{dl: &v}:
				default:
				}
			},
			func(mbps float64) {
				v := mbps
				select {
				case ch <- progressUpdate{ul: &v}:
				default:
				}
			},
		)
		if err != nil {
			return msgError{err: err}
		}
		return msgResult{result: result}
	}
}

func (m *Model) fetchISP() tea.Cmd {
	return func() tea.Msg {
		info, err := network.FetchIPInfo()
		if err != nil {
			return msgISPDone{info: nil}
		}
		return msgISPDone{info: info}
	}
}
