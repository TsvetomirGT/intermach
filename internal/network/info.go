package network

import (
	"encoding/json"
	"net/http"
	"time"
)

type IPInfo struct {
	IP   string `json:"ip"`
	Org  string `json:"org"`
	City string `json:"city"`
}

func FetchIPInfo() (*IPInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
