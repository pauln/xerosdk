package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/quickaco/xerosdk/helpers"
)

const (
	batchPaymentURL = "https://api.xero.com/api.xro/2.0/BatchPayments"
)

// BatchPayment type will keep information related with a bank
type BatchPayment struct {

	// A user defined bank account number.
	BankAccountNumber string `json:"BankAccountNumber,omitempty"`

	// Full name of bank account
	BankAccountName string `json:"BankAccountName,omitempty"`

	// Details of the Batch payment
	Details string `json:"Details,omitempty"`

	// Code of the Batch payment
	Code string `json:"Code,omitempty"`

	// Reference of the Batch payment
	Reference string `json:"Reference,omitempty"`
}

func unmarshalBatchPayment(batchPaymentResponseBytes []byte) ([]BatchPayment, error) {
	response := struct {
		Payments []BatchPayment `json:"BatchPayments,omitempty"`
	}{}
	err := json.Unmarshal(batchPaymentResponseBytes, &response)
	if err != nil {
		return nil, err
	}
	return response.Payments, nil
}

// FindBatchPayments will get all the batch payments
func FindBatchPayments(cl *http.Client) ([]BatchPayment, error) {
	batchPayments, err := helpers.Find(cl, batchPaymentURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalBatchPayment(batchPayments)
}
