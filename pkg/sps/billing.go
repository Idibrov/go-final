package sps

import (
	"io/ioutil"
	"math"
	"os"
)

type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"recurring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}


const (
	BILLING_CREATE_CUSTOMER = iota
	BILLING_PURCHASE
	BILLING_PAYOUT
	BILLING_RECURRING
	BILLING_FRAUD_CONTROL
	BILLING_CHECKOUT_PAGE
)


func GetStatusBilling(csvPath string) (BillingData, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return BillingData{}, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return BillingData{}, err
	}

	err = file.Close()
	if err != nil {
		return BillingData{}, err
	}

	return parseBilling(content), nil

}


func parseBilling(mask []byte) BillingData {
	if len(mask) != 6 {
		return BillingData{}
	}

	var bitMask int8 = 0
	for i, bit := range mask {
		if bit == '1' {
			bitMask += int8(math.Pow(2, float64(len(mask)-i-1)))
		}
	}

	createCustomer := bitMask>>BILLING_CREATE_CUSTOMER&1 == 1
	purchase := bitMask>>BILLING_PURCHASE&1 == 1
	payout := bitMask>>BILLING_PAYOUT&1 == 1
	recurring := bitMask>>BILLING_RECURRING&1 == 1
	fraudControl := bitMask>>BILLING_FRAUD_CONTROL&1 == 1
	checkoutPage := bitMask>>BILLING_CHECKOUT_PAGE&1 == 1

	return BillingData{
		CreateCustomer: createCustomer,
		Purchase:       purchase,
		Payout:         payout,
		Recurring:      recurring,
		FraudControl:   fraudControl,
		CheckoutPage:   checkoutPage,
	}
}