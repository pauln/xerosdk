package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/quickaco/xerosdk/helpers"
)

const (
	currencyURL = "https://api.xero.com/api.xro/2.0/Currencies"
)

//Currency is the local currency set up to be used in Xero
type Currency struct {

	// 3 letter alpha code for the currency – see list of currency codes
	Code string `json:"Code,omitempty"`

	// Name of Currency
	Description string `json:"Description,omitempty"`
}

//Currencies is a collection of Currencies
type Currencies struct {
	Currencies []Currency `json:"Currencies,omitempty"`
}

func unmarshalCurrencies(currencyResponseBytes []byte) (*Currencies, error) {
	var currencyResponse *Currencies
	err := json.Unmarshal(currencyResponseBytes, &currencyResponse)
	if err != nil {
		return nil, err
	}

	return currencyResponse, err
}

// FindCurrencies will get all currencies
func FindCurrencies(cl *http.Client) (*Currencies, error) {
	currencyBytes, err := helpers.Find(cl, currencyURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return unmarshalCurrencies(currencyBytes)
}

// Create will create a new currency on Xero
func (c *Currencies) Create(cl *http.Client) (*Currencies, error) {
	buf, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	currencyBytes, err := helpers.Create(cl, currencyURL, buf)
	if err != nil {
		return nil, err
	}
	return unmarshalCurrencies(currencyBytes)
}
