package sps

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}


const (
	COUNTRY_SMS = iota
	BANDWIDTH_SMS
	RESPONSE_TIME_SMS
	PROVIDER_SMS
)


func GetStatusSMS(csvPath string) ([]SMSData, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return []SMSData{}, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return []SMSData{}, err
	}

	err = file.Close()
	if err != nil {
		return []SMSData{}, err
	}

	reader := strings.NewReader(string(content))
	scanner := bufio.NewScanner(reader)
	SMSList := make([]SMSData, 0)

	for scanner.Scan() {
		line := scanner.Text()
		sms, ok := parseSMS(line)

		if ok {
			SMSList = append(SMSList, sms)
		}
	}
	return SMSList, nil
}


func parseSMS(line string) (SMSData, bool) {
	sms := strings.Split(line, ";")

	switch {
	case len(sms) < 4:
		fallthrough
	case !isValidCountry(sms[COUNTRY_SMS]):
		fallthrough
	case !isValidBandwidth(sms[BANDWIDTH_SMS]):
		fallthrough
	case !isValidResponseTime(sms[RESPONSE_TIME_SMS]):
		fallthrough
	case !isValidSMSProvider(sms[PROVIDER_SMS]):
		return SMSData{}, false
	}

	return SMSData{
		Country:      sms[COUNTRY_SMS],
		Bandwidth:    sms[BANDWIDTH_SMS],
		ResponseTime: sms[RESPONSE_TIME_SMS],
		Provider:     sms[PROVIDER_SMS],
	}, true
}