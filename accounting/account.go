package accounting

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/quickaco/xerosdk/helpers"
)

const (
	accountsURL = "https://api.xero.com/api.xro/2.0/Accounts"
)

//Account represents individual accounts in a Xero organisation
type Account struct {

	// Customer defined alpha numeric account code e.g 200 or SALES (max length = 10)
	Code string `json:"Code,omitempty" xml:"Code,omitempty"`

	// Name of account (max length = 150)
	Name string `json:"Name,omitempty" xml:"Name,omitempty"`

	// See Account Types
	Type string `json:"Type,omitempty" xml:"Type,omitempty"`

	// For bank accounts only (Account Type BANK)
	BankAccountNumber string `json:"BankAccountNumber,omitempty" xml:"BankAccountNumber,omitempty"`

	// Accounts with a status of ACTIVE can be updated to ARCHIVED. See Account Status Codes
	Status string `json:"Status,omitempty" xml:"Status,omitempty"`

	// Description of the Account. Valid for all types of accounts except bank accounts (max length = 4000)
	Description string `json:"Description,omitempty" xml:"Description,omitempty"`

	// For bank accounts only. See Bank Account types
	BankAccountType string `json:"BankAccountType,omitempty" xml:"BankAccountType,omitempty"`

	// For bank accounts only
	CurrencyCode string `json:"CurrencyCode,omitempty" xml:"CurrencyCode,omitempty"`

	// See Tax Types
	TaxType string `json:"TaxType,omitempty" xml:"TaxType,omitempty"`

	// Boolean – describes whether account can have payments applied to it
	EnablePaymentsToAccount bool `json:"EnablePaymentsToAccount,omitempty" xml:"EnablePaymentsToAccount,omitempty"`

	// Boolean – describes whether account code is available for use with expense claims
	ShowInExpenseClaims bool `json:"ShowInExpenseClaims,omitempty" xml:"ShowInExpenseClaims,omitempty"`

	// The Xero identifier for an account – specified as a string following the endpoint name e.g. /297c2dc5-cc47-4afd-8ec8-74990b8761e9
	AccountID string `json:"AccountID,omitempty" xml:"AccountID,omitempty"`

	// See Account Class Types
	Class string `json:"Class,omitempty" xml:"-"`

	// If this is a system account then this element is returned. See System Account types. Note that non-system accounts may have this element set as either “” or null.
	SystemAccount string `json:"SystemAccount,omitempty" xml:"-"`

	// Shown if set
	ReportingCode string `json:"ReportingCode,omitempty" xml:"-"`

	// Shown if set
	ReportingCodeName string `json:"ReportingCodeName,omitempty" xml:"-"`

	// boolean to indicate if an account has an attachment (read only)
	HasAttachments bool `json:"HasAttachments,omitempty" xml:"-"`

	// Last modified date UTC format
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty" xml:"-"`
}

//Accounts contains a collection of Accounts
type Accounts struct {
	Accounts []Account `json:"Accounts,omitempty" xml:"Account,omitempty"`
}

//The Xero API returns Dates based on the .Net JSON date format available at the time of development
//We need to convert these to a more usable format - RFC3339 for consistency with what the API expects to recieve
func (a *Accounts) convertDates() error {
	var err error
	for n := len(a.Accounts) - 1; n >= 0; n-- {
		a.Accounts[n].UpdatedDateUTC, err = helpers.DotNetJSONTimeToRFC3339(a.Accounts[n].UpdatedDateUTC, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalAccount(accountResponseBytes []byte) (*Accounts, error) {
	var accountResponse *Accounts
	err := json.Unmarshal(accountResponseBytes, &accountResponse)
	if err != nil {
		return nil, err
	}

	err = accountResponse.convertDates()
	if err != nil {
		return nil, err
	}

	return accountResponse, err
}

// FindAccounts will retrieve all the accounts
func FindAccounts(cl *http.Client) (*Accounts, error) {
	request, err := http.NewRequest(http.MethodGet, accountsURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	accountResponseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return unmarshalAccount(accountResponseBytes)
}
