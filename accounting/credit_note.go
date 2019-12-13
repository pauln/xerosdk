package accounting

//CreditNote an be raised directly against a customer or supplier,
//allowing the customer or supplier to be held in credit until a future invoice or bill is raised
type CreditNote struct {

	// See Credit Note Types
	Type string `json:"Type,omitempty" xml:"Type,omitempty"`

	// See Contacts
	Contact Contact `json:"Contact" xml:"Contact"`

	// The date the credit note is issued YYYY-MM-DD.
	// If the Date element is not specified then it will default
	// to the current date based on the timezone setting of the organisation
	Date string `json:"DateString,omitempty" xml:"Date,omitempty"`

	// See Credit Note Status Codes
	Status string `json:"Status,omitempty" xml:"Status,omitempty"`

	// See Invoice Line Amount Types
	LineAmountTypes string `json:"LineAmountTypes,omitempty" xml:"LineAmountTypes,omitempty"`

	// See Invoice Line Items
	LineItems []LineItem `json:"LineItems,omitempty" xml:"LineItems>LineItem,omitempty"`

	// The subtotal of the credit note excluding taxes
	SubTotal float64 `json:"SubTotal,omitempty" xml:"SubTotal,omitempty"`

	// The total tax on the credit note
	TotalTax float64 `json:"TotalTax,omitempty" xml:"TotalTax,omitempty"`

	// The total of the Credit Note(subtotal + total tax)
	Total float64 `json:"Total,omitempty" xml:"Total,omitempty"`

	// UTC timestamp of last update to the credit note
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty" xml:"-"`

	// Currency used for the Credit Note
	CurrencyCode string `json:"CurrencyCode,omitempty" xml:"CurrencyCode,omitempty"`

	// Date when credit note was fully paid(UTC format)
	FullyPaidOnDate string `json:"FullyPaidOnDate,omitempty" xml:"-"`

	// Xero generated unique identifier
	CreditNoteID string `json:"CreditNoteID,omitempty" xml:"CreditNoteID,omitempty"`

	// ACCRECCREDIT – Unique alpha numeric code identifying credit note (when missing will auto-generate from your Organisation Invoice Settings)
	CreditNoteNumber string `json:"CreditNoteNumber,omitempty" xml:"CreditNoteNumber,omitempty"`

	// ACCRECCREDIT only – additional reference number
	Reference string `json:"Reference,omitempty" xml:"Reference,omitempty"`

	// boolean to indicate if a credit note has been sent to a contact via the Xero app (currently read only)
	SentToContact bool `json:"SentToContact,omitempty" xml:"SentToContact,omitempty"`

	// The currency rate for a multicurrency invoice. If no rate is specified, the XE.com day rate is used
	CurrencyRate float64 `json:"CurrencyRate,omitempty" xml:"CurrencyRate,omitempty"`

	// The remaining credit balance on the Credit Note
	RemainingCredit float64 `json:"RemainingCredit,omitempty" xml:"-"`

	// See Allocations
	Allocations *[]Allocation `json:"Allocations,omitempty" xml:"-"`

	// See BrandingThemes
	BrandingThemeID string `json:"BrandingThemeID,omitempty" xml:"BrandingThemeID,omitempty"`

	// boolean to indicate if a credit note has an attachment
	HasAttachments bool `json:"HasAttachments,omitempty" xml:"-"`
}

//CreditNotes is a collection of CreditNote
type CreditNotes struct {
	CreditNotes []CreditNote `json:"CreditNotes" xml:"CreditNote"`
}
