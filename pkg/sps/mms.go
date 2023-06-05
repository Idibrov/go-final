package sps

import (
	"encoding/json"
	"io"
	"net/http"
)

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}


func MMSRequest(addr string) ([]MMSData, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return []MMSData{}, err
	}

	if resp.StatusCode != 200 {
		return []MMSData{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []MMSData{}, err
	}

	data := make([]MMSData, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return []MMSData{}, err
	}

	filteredData := make([]MMSData, 0)
	for _, mms := range data {
		switch {
		case !isValidCountry(mms.Country):
			fallthrough
		case !isValidBandwidth(mms.Bandwidth):
			fallthrough
		case !isValidResponseTime(mms.ResponseTime):
			fallthrough
		case !isValidSMSProvider(mms.Provider):
		default:
			filteredData = append(filteredData, mms)
		}
	}

	return filteredData, nil
}