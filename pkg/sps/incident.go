package sps

import (
	"encoding/json"
	"io"
	"net/http"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"` 
}


func IncidentRequest(addr string) ([]IncidentData, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return []IncidentData{}, err
	}

	if resp.StatusCode != 200 {
		return []IncidentData{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []IncidentData{}, err
	}

	data := make([]IncidentData, 0)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return []IncidentData{}, err
	}

	filteredData := make([]IncidentData, 0)
	for _, incident := range data {
		switch {
		case !isValidStatus(incident.Status):
		default:
			filteredData = append(filteredData, incident)
		}
	}

	return filteredData, nil
}