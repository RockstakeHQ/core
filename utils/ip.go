package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type IPInfo struct {
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
	Query    string `json:"query"`
}

func GetClientTimezone() (string, error) {
	response, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var ipInfo IPInfo
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return "", err
	}
	return ipInfo.Timezone, nil
}
