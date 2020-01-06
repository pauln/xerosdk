package accounting

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/helpers"
)

const (
	creditNotesURL = "https://api.xero.com/api.xro/2.0/CreditNotes"
)

//CreditNote an be raised directly against a customer or supplier,
//allowing the customer or supplier to be held in credit until a future invoice or bill is raised
type CreditNote struct {

	// See Credit Note Types
	Type string `json:"Type,omitempty"`

	// See Contacts
	Contact Contact `json:"Contact"`

	// The date the credit note is issued YYYY-MM-DD.
	// If the Date element is not specified then it will default
	// to the current date based on the timezone setting of the organisation
	Date string `json:"DateString,omitempty"`

	// See Credit Note Status Codes
	Status string `json:"Status,omitempty"`

	// See Invoice Line Amount Types
	LineAmountTypes string `json:"LineAmountTypes,omitempty"`

	// See Invoice Line Items
	LineItems []LineItem `json:"LineItems,omitempty"`

	// The subtotal of the credit note excluding taxes
	SubTotal float64 `json:"SubTotal,omitempty"`

	// The total tax on the credit note
	TotalTax float64 `json:"TotalTax,omitempty"`

	// The total of the Credit Note(subtotal + total tax)
	Total float64 `json:"Total,omitempty"`

	// UTC timestamp of last update to the credit note
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty"`

	// Currency used for the Credit Note
	CurrencyCode string `json:"CurrencyCode,omitempty"`

	// Date when credit note was fully paid(UTC format)
	FullyPaidOnDate string `json:"FullyPaidOnDate,omitempty"`

	// Xero generated unique identifier
	CreditNoteID string `json:"CreditNoteID,omitempty"`

	// ACCRECCREDIT – Unique alpha numeric code identifying credit note (when missing will auto-generate from your Organisation Invoice Settings)
	CreditNoteNumber string `json:"CreditNoteNumber,omitempty"`

	// ACCRECCREDIT only – additional reference number
	Reference string `json:"Reference,omitempty"`

	// boolean to indicate if a credit note has been sent to a contact via the Xero app (currently read only)
	SentToContact bool `json:"SentToContact,omitempty"`

	// The currency rate for a multicurrency invoice. If no rate is specified, the XE.com day rate is used
	CurrencyRate float64 `json:"CurrencyRate,omitempty"`

	// The remaining credit balance on the Credit Note
	RemainingCredit float64 `json:"RemainingCredit,omitempty"`

	// See Allocations
	Allocations *[]Allocation `json:"Allocations,omitempty"`

	// See BrandingThemes
	BrandingThemeID string `json:"BrandingThemeID,omitempty"`

	// boolean to indicate if a credit note has an attachment
	HasAttachments bool `json:"HasAttachments,omitempty"`
}

//CreditNotes is a collection of CreditNote
type CreditNotes struct {
	CreditNotes []CreditNote `json:"CreditNotes"`
}

//The Xero API returns Dates based on the .Net JSON date format available at the time of development
//We need to convert these to a more usable format - RFC3339 for consistency with what the API expects to recieve
func (c *CreditNotes) convertDates() error {
	var err error
	for n := len(c.CreditNotes) - 1; n >= 0; n-- {
		c.CreditNotes[n].UpdatedDateUTC, err = helpers.DotNetJSONTimeToRFC3339(c.CreditNotes[n].UpdatedDateUTC, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalCreditNote(creditNoteResponseBytes []byte) (*CreditNotes, error) {
	var creditNoteResponse *CreditNotes
	err := json.Unmarshal(creditNoteResponseBytes, &creditNoteResponse)
	if err != nil {
		return nil, err
	}

	err = creditNoteResponse.convertDates()
	if err != nil {
		return nil, err
	}

	return creditNoteResponse, err
}

// Create will create accounts given an Accounts struct
func (c *CreditNotes) Create(cl *http.Client) (*CreditNotes, error) {
	buf, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	creditNotesBytes, err := helpers.Create(cl, creditNotesURL, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalCreditNote(creditNotesBytes)
}

// Update will update an account given an Accounts struct
// This will only handle single account - you cannot update multiple accounts in a single call
func (c *CreditNote) Update(cl *http.Client) (*CreditNotes, error) {
	cn := CreditNotes{
		CreditNotes: []CreditNote{*c},
	}
	buf, err := json.Marshal(cn)
	if err != nil {
		return nil, err
	}
	creditNotesBytes, err := helpers.Update(cl, creditNotesURL+"/"+c.CreditNoteID, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalCreditNote(creditNotesBytes)
}

// FindCreditNotes will get all CreditNotes. These Credit Notes will not have details like line items by default.
// If you need details then then add a 'page' querystringParameter and get 100 Credit Notes at a time
// additional querystringParameters such as where, page, order can be added as a map
func FindCreditNotes(cl *http.Client, queryParameters map[string]string) (*CreditNotes, error) {
	creditNotes, err := helpers.Find(cl, creditNotesURL, nil, queryParameters)
	if err != nil {
		return nil, err
	}

	return unmarshalCreditNote(creditNotes)
}

// FindCreditNotesModifiedSince will get all Credit Notes modified after a specified date.
// These Credit Notes will not have details like line items by default.
// If you need details then then add a 'page' querystringParameter and get 100 Credit Notes at a time
// additional querystringParameters such as where, page, order can be added as a map
func FindCreditNotesModifiedSince(cl *http.Client, modifiedSince time.Time, queryParameters map[string]string) (*CreditNotes, error) {
	additionalHeaders := map[string]string{}
	additionalHeaders["If-Modified-Since"] = modifiedSince.Format(time.RFC3339)

	creditNotes, err := helpers.Find(cl, creditNotesURL, additionalHeaders, queryParameters)
	if err != nil {
		return nil, err
	}

	return unmarshalCreditNote(creditNotes)
}

// FindCreditNote will get a single creditNote - creditNoteID can be a GUID for a creditNote or a creditNote number
func FindCreditNote(cl *http.Client, creditNoteID uuid.UUID) (*CreditNote, error) {
	creditNotes, err := helpers.Find(cl, creditNotesURL+"/"+creditNoteID.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	notes, err := unmarshalCreditNote(creditNotes)
	if err != nil {
		return nil, err
	}
	return &notes.CreditNotes[0], nil
}
