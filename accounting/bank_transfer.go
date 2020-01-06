package accounting

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/helpers"
)

const (
	bankTransferURL = "https://api.xero.com/api.xro/2.0/BankTransfers"
)

//BankTransfer is a record of monies transferred from one bank account to another
type BankTransfer struct {

	//
	Amount float64 `json:"Amount" xml:"Amount"`

	// The date of the Transfer YYYY-MM-DD
	Date string `json:"Date,omitempty" xml:"Date,omitempty"`

	// The identifier of the Bank Transfer
	BankTransferID string `json:"BankTransferID,omitempty" xml:"BankTransferID,omitempty"`

	// The currency rate
	CurrencyRate float64 `json:"CurrencyRate,omitempty" xml:"CurrencyRate,omitempty"`

	// The Bank Transaction ID for the source account
	FromBankTransactionID string `json:"FromBankTransactionID,omitempty" xml:"FromBankTransactionID,omitempty"`

	// The Bank Transaction ID for the destination account
	ToBankTransactionID string `json:"ToBankTransactionID,omitempty" xml:"ToBankTransactionID,omitempty"`

	// Boolean to indicate if a Bank Transfer has an attachment
	HasAttachments bool `json:"HasAttachments,omitempty" xml:"HasAttachments,omitempty"`

	// UTC timestamp of creation date of bank transfer
	CreatedDateUTC string `json:"CreatedDateUTC,omitempty" xml:"CreatedDateUTC,omitempty"`

	// The source BankAccount
	FromBankAccount BankAccount `json:"FromBankAccount,omitempty" xml:"FromBankAccount,omitempty"`

	// The destination BankAccount
	ToBankAccount BankAccount `json:"ToBankAccount,omitempty" xml:"ToBankAccount,omitempty"`
}

//BankTransfers contains a collection of BankTransfers
type BankTransfers struct {
	BankTransfers []BankTransfer `json:"BankTransfers" xml:"BankTransfer"`
}

//The Xero API returns Dates based on the .Net JSON date format available at the time of development
//We need to convert these to a more usable format - RFC3339 for consistency with what the API expects to recieve
func (b *BankTransfers) convertDates() error {
	var err error
	for n := len(b.BankTransfers) - 1; n >= 0; n-- {
		b.BankTransfers[n].Date, err = helpers.DotNetJSONTimeToRFC3339(b.BankTransfers[n].Date, false)
		if err != nil {
			return err
		}
		b.BankTransfers[n].CreatedDateUTC, err = helpers.DotNetJSONTimeToRFC3339(b.BankTransfers[n].CreatedDateUTC, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalBankTransfer(bankTransferResponseBytes []byte) (*BankTransfers, error) {
	var bankTransferResponse *BankTransfers
	err := json.Unmarshal(bankTransferResponseBytes, &bankTransferResponse)
	if err != nil {
		return nil, err
	}

	err = bankTransferResponse.convertDates()
	if err != nil {
		return nil, err
	}

	return bankTransferResponse, err
}

// FindBankTransfersModifiedSince will get all BankTransfers modified after a specified date.
// These BankTransfers will not have details like default line items by default.
// If you need details then add a 'page' querystringParameter and get 100 BankTransfers at a time
// additional querystringParameters such as where and order can be added as a map
func FindBankTransfersModifiedSince(cl *http.Client, modifiedSince time.Time, queryParameters map[string]string) (*BankTransfers, error) {
	additionalHeaders := map[string]string{}
	additionalHeaders["If-Modified-Since"] = modifiedSince.Format(time.RFC3339)

	bankTransferBytes, err := helpers.Find(cl, bankTransferURL, additionalHeaders, queryParameters)
	if err != nil {
		return nil, err
	}

	return unmarshalBankTransfer(bankTransferBytes)
}

// FindBankTransfers will get all BankTransfers. These BankTransfer will not have details like line items by default.
// If you need details then add a 'page' querystringParameter and get 100 BankTransfers at a time
// additional querystringParameters such as where and order can be added as a map
func FindBankTransfers(cl *http.Client, queryParameters map[string]string) (*BankTransfers, error) {
	bankTransferBytes, err := helpers.Find(cl, bankTransferURL, nil, queryParameters)
	if err != nil {
		return nil, err
	}

	return unmarshalBankTransfer(bankTransferBytes)
}

// FindBankTransfer will get a single bankTransfer - bankTransferID can be a GUID for an bankTransfer or an bankTransfer number
func FindBankTransfer(cl *http.Client, bankTransferID uuid.UUID) (*BankTransfer, error) {
	bankTransferBytes, err := helpers.Find(cl, bankTransferURL+"/"+bankTransferID.String(), nil, nil)
	if err != nil {
		return nil, err
	}
	b, err := unmarshalBankTransfer(bankTransferBytes)
	if err != nil {
		return nil, err
	}
	if len(b.BankTransfers) > 0 {
		return &b.BankTransfers[0], nil
	}
	return nil, nil
}

// Create will create bankTransfers given a BankTransfers struct
func (b *BankTransfers) Create(cl *http.Client) (*BankTransfers, error) {
	buf, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	bankTransferBytes, err := helpers.Create(cl, bankTransferURL, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalBankTransfer(bankTransferBytes)
}
