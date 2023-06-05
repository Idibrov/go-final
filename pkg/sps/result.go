package sps

import (
	"sort"
)

type ResultT struct {
	Status bool       `json:"status"`
	Data   ResultSetT `json:"data"`
	Error  string     `json:"error"`
}

type ResultSetT struct {
	SMS       [][]SMSData              `json:"sms"`
	MMS       [][]MMSData              `json:"mms"`
	VoiceCall []VoiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []IncidentData           `json:"incident"`
}


var (
	smsFilename     = "/sms.data"
	emailFilename   = "/email.data"
	billingFilename = "/billing.data"
	voiceFilename   = "/voice.data"
)


var (
	mmsUrl       = "/mms"
	supportUrl   = "/support"
	accendentUrl = "/accendent"
)


const (
	LOW_LOAD  = 1
	AVG_LOAD  = 2
	HIGH_LOAD = 3
)


func GetResultData(apiAddr, dataPath string) (ResultSetT, error) {
	var result ResultSetT

	var err error
	result.SMS, err = getResultSMS(dataPath + smsFilename)
	if err != nil {
		return ResultSetT{}, err
	}

	result.MMS, err = getResultMMS(apiAddr + mmsUrl)
	if err != nil {
		return ResultSetT{}, err
	}

	result.VoiceCall, err = getResultVoice(dataPath + voiceFilename)
	if err != nil {
		return ResultSetT{}, err
	}

	result.Email, err = getResultEmail(dataPath + emailFilename)
	if err != nil {
		return ResultSetT{}, err
	}

	result.Billing, err = getResultBilling(dataPath + billingFilename)
	if err != nil {
		return ResultSetT{}, err
	}

	result.Support, err = getResultSupport(apiAddr + supportUrl)
	if err != nil {
		return ResultSetT{}, err
	}

	result.Incidents, err = getResultIncidents(apiAddr + accendentUrl)
	if err != nil {
		return ResultSetT{}, err
	}

	return result, nil
}


func getResultSMS(smsPath string) ([][]SMSData, error) {
	data, err := GetStatusSMS(smsPath)
	if err != nil {
		return [][]SMSData{}, err
	}

	for i := range data {
		data[i].Country = codeToCountry(data[i].Country)
	}

	smsCountrySort := data
	smsProviderSort := make([]SMSData, len(data))
	copy(smsProviderSort, smsCountrySort)

	sort.SliceStable(smsCountrySort, func(i, j int) bool {
		return smsCountrySort[i].Country < smsCountrySort[j].Country
	})
	sort.SliceStable(smsProviderSort, func(i, j int) bool {
		return smsProviderSort[i].Provider < smsProviderSort[j].Provider
	})

	return [][]SMSData{smsProviderSort, smsCountrySort}, nil
}


func getResultMMS(addr string) ([][]MMSData, error) {
	data, err := MMSRequest(addr)
	if err != nil {
		return [][]MMSData{}, err
	}

	for i := range data {
		data[i].Country = codeToCountry(data[i].Country)
	}

	mmsCountrySort := data
	mmsProviderSort := make([]MMSData, len(data))
	copy(mmsProviderSort, mmsCountrySort)

	sort.SliceStable(mmsCountrySort, func(a, b int) bool {
		return mmsCountrySort[a].Country < mmsCountrySort[b].Country
	})
	sort.SliceStable(mmsProviderSort, func(a, b int) bool {
		return mmsProviderSort[a].Provider < mmsProviderSort[b].Provider
	})

	return [][]MMSData{mmsProviderSort, mmsCountrySort}, nil
}


func getResultVoice(voicePath string) ([]VoiceCallData, error) {
	voice, err := GetStatusVoice(voicePath)
	if err != nil {
		return []VoiceCallData{}, err
	}
	return voice, nil
}


func getResultEmail(emailPath string) (map[string][][]EmailData, error) {
	data, err := GetStatusEmail(emailPath)
	if err != nil {
		return map[string][][]EmailData{}, err
	}
	resultEmail := make(map[string][][]EmailData, 0)

	for code := range countryCode {
		slowProviders, fastProviders := GetSlowFastEmailProvider(data, code)
		if len(slowProviders) > 2 && len(fastProviders) > 2 {
			resultEmail[code] = [][]EmailData{slowProviders, fastProviders}
		}
	}

	return resultEmail, nil
}


func getResultBilling(billingPath string) (BillingData, error) {
	billing, err := GetStatusBilling(billingPath)
	if err != nil {
		return BillingData{}, err
	}

	return billing, nil
}


func getResultSupport(addr string) ([]int, error) {
	data, err := SupportRequest(addr)
	if err != nil {
		return []int{}, err
	}

	if len(data) == 0 {
		return []int{}, err
	}

	activeTickets := 0
	for _, support := range data {
		activeTickets += support.ActiveTickets
	}

	avgTime := activeTickets * 60 / 18
	load := LOW_LOAD
	switch {
	case activeTickets >= 9 && activeTickets <= 16:
		load = AVG_LOAD
	case activeTickets > 16:
		load = HIGH_LOAD
	}

	return []int{load, avgTime}, nil
}


func getResultIncidents(addr string) ([]IncidentData, error) {
	data, err := IncidentRequest(addr)
	if err != nil {
		return []IncidentData{}, nil
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Status < data[j].Status
	})

	return data, nil
}