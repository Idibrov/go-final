package sps

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type VoiceCallData struct {
	Country        string  `json:"country"`
	Load           string  `json:"bandwidth"`
	ResponseTime   string  `json:"response_time"`
	Provider       string  `json:"provider"`
	Stability      float32 `json:"connection_stability"`
	TTFB           int     `json:"ttfb"`
	Purity         int     `json:"voice_purity"`
	MedianDuration int     `json:"median_of_call_time"`
}


const (
	COUNTRY_VOICE = iota
	LOAD_VOICE
	RESPONSE_TIME_VOICE
	PROVIDER_VOICE
	STABILITY_VOICE
	TTFB_VOICE
	PURITY_VOICE
	MEDIAN_DURATION_VOICE
)


func GetStatusVoice(csvPath string) ([]VoiceCallData, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return []VoiceCallData{}, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return []VoiceCallData{}, err
	}

	reader := strings.NewReader(string(content))
	scanner := bufio.NewScanner(reader)
	VoiceList := make([]VoiceCallData, 0)

	for scanner.Scan() {
		line := scanner.Text()
		voice, ok := parseVoiceData(line)

		if ok {
			VoiceList = append(VoiceList, voice)
		}
	}

	return VoiceList, nil
}


func parseVoiceData(line string) (VoiceCallData, bool) {
	voice := strings.Split(line, ";")

	switch {
	case len(voice) != 8:
		fallthrough
	case !isValidCountry(voice[COUNTRY_VOICE]):
		fallthrough
	case !isValidLoad(voice[LOAD_VOICE]):
		fallthrough
	case !isValidResponseTime(voice[RESPONSE_TIME_VOICE]):
		fallthrough
	case !isValidVoiceProvider(voice[PROVIDER_VOICE]):
		fallthrough
	case !isValidStability(voice[STABILITY_VOICE]):
		fallthrough
	case !isValidPurity(voice[PURITY_VOICE]):
		fallthrough
	case !isValidTTFB(voice[TTFB_VOICE]):
		fallthrough
	case !isMedianDuration(voice[MEDIAN_DURATION_VOICE]):
		return VoiceCallData{}, false
	}

	load := voice[LOAD_VOICE]
	responseTime := voice[RESPONSE_TIME_VOICE]
	stability64, _ := strconv.ParseFloat(voice[STABILITY_VOICE], 32)
	ttfb, _ := strconv.Atoi(voice[RESPONSE_TIME_VOICE])
	purity, _ := strconv.Atoi(voice[PURITY_VOICE])
	medianDuration, _ := strconv.Atoi(voice[MEDIAN_DURATION_VOICE])

	return VoiceCallData{
		Country:        voice[COUNTRY_VOICE],
		Load:           load,
		ResponseTime:   responseTime,
		Provider:       voice[PROVIDER_VOICE],
		Stability:      float32(stability64),
		TTFB:           ttfb,
		Purity:         purity,
		MedianDuration: medianDuration,
	}, true
}