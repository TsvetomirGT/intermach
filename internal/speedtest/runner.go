package speedtest

import (
	"github.com/showwin/speedtest-go/speedtest"
)

type Result struct {
	DownloadMbps  float64
	UploadMbps    float64
	LatencyMs     float64
	ServerName    string
	ServerCountry string
}

type ProgressUpdate struct {
	Mbps float64
}

// Run performs the full speed test synchronously, calling the provided callbacks
// for live progress. Designed to be called from a goroutine inside a tea.Cmd.
func Run(
	onDownloadProgress func(mbps float64),
	onUploadProgress func(mbps float64),
) (*Result, error) {
	client := speedtest.New()

	servers, err := client.FetchServers()
	if err != nil {
		return nil, err
	}

	targets, err := servers.FindServer(nil)
	if err != nil {
		return nil, err
	}

	s := targets[0]

	// Ping
	if err := s.PingTest(nil); err != nil {
		return nil, err
	}

	// Download with live callback
	s.Context.SetCallbackDownload(func(rate speedtest.ByteRate) {
		if onDownloadProgress != nil {
			onDownloadProgress(rate.Mbps())
		}
	})
	if err := s.DownloadTest(); err != nil {
		return nil, err
	}

	// Upload with live callback
	s.Context.SetCallbackUpload(func(rate speedtest.ByteRate) {
		if onUploadProgress != nil {
			onUploadProgress(rate.Mbps())
		}
	})
	if err := s.UploadTest(); err != nil {
		return nil, err
	}

	return &Result{
		DownloadMbps:  s.DLSpeed.Mbps(),
		UploadMbps:    s.ULSpeed.Mbps(),
		LatencyMs:     float64(s.Latency.Milliseconds()),
		ServerName:    s.Name,
		ServerCountry: s.Country,
	}, nil
}
