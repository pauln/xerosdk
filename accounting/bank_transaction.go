package accounting

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/helpers"
)

const (
	bankTransactionURL = "https://api.xero.com/api.xro/2.0/BankTransactions"
)

//BankTransaction is a bank transaction
type BankTransaction struct {

	// See Bank Transaction Types
	Type string `json:"Type" xml:"Type"`

	// See Contacts
	Contact Contact `json:"Contact" xml:"Contact"`

	// See LineItems
	LineItems []LineItem `json:"LineItems" xml:"LineItems>LineItem"`

	// Boolean to show if transaction is reconciled
	IsReconciled bool `json:"IsReconciled,omitempty" xml:"IsReconciled,omitempty"`

	// Date of transaction – YYYY-MM-DD
	Date string `json:"DateString,omitempty" xml:"Date,omitempty"`

	// Reference for the transaction. Only supported for SPEND and RECEIVE transactions.
	Reference string `json:"Reference,omitempty" xml:"Reference,omitempty"`

	// The currency that bank transaction has been raised in (see Currencies). Setting currency is only supported on overpayments.
	CurrencyCode string `json:"CurrencyCode,omitempty" xml:"CurrencyCode,omitempty"`

	// Exchange rate to base currency when money is spent or received. e.g. 0.7500 Only used for bank transactions in non base currency. If this isn’t specified for non base currency accounts then either the user-defined rate (preference) or the XE.com day rate will be used. Setting currency is only supported on overpayments.
	CurrencyRate float64 `json:"CurrencyRate,omitempty" xml:"CurrencyRate,omitempty"`

	// URL link to a source document – shown as “Go to App Name”
	URL string `json:"Url,omitempty" xml:"Url,omitempty"`

	// See Bank Transaction Status Codes
	Status string `json:"Status,omitempty" xml:"Status,omitempty"`

	// Line amounts are exclusive of tax by default if you don’t specify this element. See Line Amount Types
	LineAmountTypes string `json:"LineAmountTypes,omitempty" xml:"LineAmountTypes,omitempty"`

	// Total of bank transaction excluding taxes
	SubTotal float64 `json:"SubTotal,omitempty" xml:"SubTotal,omitempty"`

	// Total tax on bank transaction
	TotalTax float64 `json:"TotalTax,omitempty" xml:"TotalTax,omitempty"`

	// Total of bank transaction tax inclusive
	Total float64 `json:"Total,omitempty" xml:"Total,omitempty"`

	// Xero generated unique identifier for bank transaction
	BankTransactionID string `json:"BankTransactionID,omitempty" xml:"BankTransactionID,omitempty"`

	// Xero Bank Account
	BankAccount BankAccount `json:"BankAccount,omitempty" xml:"BankAccount,omitempty"`

	// Xero generated unique identifier for a Prepayment. This will be returned on BankTransactions with a Type of SPEND-PREPAYMENT or RECEIVE-PREPAYMENT
	PrepaymentID string `json:"PrepaymentID,omitempty" xml:"-"`

	// Xero generated unique identifier for an Overpayment. This will be returned on BankTransactions with a Type of SPEND-OVERPAYMENT or RECEIVE-OVERPAYMENT
	OverpaymentID string `json:"OverpaymentID,omitempty" xml:"-"`

	// Last modified date UTC format
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty" xml:"-"`

	// Boolean to indicate if a bank transaction has an attachment
	HasAttachments bool `json:"HasAttachments,omitempty" xml:"-"`
}

//BankTransactions contains a collection of BankTransactions
type BankTransactions struct {
	BankTransactions []BankTransaction `json:"BankTransactions" xml:"BankTransaction"`
}

//The Xero API returns Dates based on the .Net JSON date format available at the time of development
//We need to convert these to a more usable format - RFC3339 for consistency with what the API expects to recieve
func (b *BankTransactions) convertDates() error {
	var err error
	for n := len(b.BankTransactions) - 1; n >= 0; n-- {
		b.BankTransactions[n].UpdatedDateUTC, err = helpers.DotNetJSONTimeToRFC3339(b.BankTransactions[n].UpdatedDateUTC, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalBankTransaction(bankTransactionResponseBytes []byte) (*BankTransactions, error) {
	var bankTransactionResponse *BankTransactions
	err := json.Unmarshal(bankTransactionResponseBytes, &bankTransactionResponse)
	if err != nil {
		return nil, err
	}

	err = bankTransactionResponse.convertDates()
	if err != nil {
		return nil, err
	}

	return bankTransactionResponse, err
}

// FindBankTransactions will get all BankTransactions. These BankTransaction will not have details like line items by default.
// If you need details then then add a 'page' querystringParameter and get 100 BankTransactions at a time
// additional querystringParameters such as where, page, order can be added as a map
func FindBankTransactions(cl *http.Client, queryParameters map[string]string) (*BankTransactions, error) {
	bankTransactionsBytes, err := helpers.Find(cl, bankTransactionURL, nil, queryParameters)
	if err != nil {
		return nil, err
	}

	return unmarshalBankTransaction(bankTransactionsBytes)
}

// FindBankTransactionsModifiedSince will get all BankTransactions modified after a specified date.
// These BankTransactions will not have details like default account codes and tracking categories by default.
// If you need details then then add a 'page' querystringParameter and get 100 BankTransactions at a time
// additional querystringParameters such as where, page, order can be added as a map
func FindBankTransactionsModifiedSince(cl *http.Client, modifiedSince time.Time, queryParameters map[string]string) (*BankTransactions, error) {
	additionalHeaders := map[string]string{}
	additionalHeaders["If-Modified-Since"] = modifiedSince.Format(time.RFC3339)

	bankTransactionsBytes, err := helpers.Find(cl, bankTransactionURL, additionalHeaders, queryParameters)
	if err != nil {
		return nil, err
	}

	return unmarshalBankTransaction(bankTransactionsBytes)
}

//FindBankTransaction will get a single BankTransaction - BankTransactionID can be a GUID for an BankTransaction or an BankTransaction number
func FindBankTransaction(cl *http.Client, bankTransactionID uuid.UUID) (*BankTransaction, error) {
	bankTransactionBytes, err := helpers.Find(cl, bankTransactionURL+"/"+bankTransactionID.String(), nil, nil)
	if err != nil {
		return nil, err
	}
	b, err := unmarshalBankTransaction(bankTransactionBytes)
	if err != nil {
		return nil, err
	}
	if len(b.BankTransactions) > 0 {
		return &b.BankTransactions[0], nil
	}
	return nil, nil
}

// Create will create accounts given an Accounts struct
func (b *BankTransactions) Create(cl *http.Client) (*BankTransactions, error) {
	buf, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	bankTransactionBytes, err := helpers.Create(cl, bankTransactionURL, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalBankTransaction(bankTransactionBytes)
}

// Update will update an account given an Accounts struct
// This will only handle single account - you cannot update multiple accounts in a single call
func (b *BankTransaction) Update(cl *http.Client) (*BankTransactions, error) {
	bt := BankTransactions{
		BankTransactions: []BankTransaction{*b},
	}
	buf, err := json.Marshal(bt)
	if err != nil {
		return nil, err
	}
	bankTransactionBytes, err := helpers.Update(cl, bankTransactionURL+"/"+b.BankTransactionID, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalBankTransaction(bankTransactionBytes)
}
