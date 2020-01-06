package accounting

//Prepayment are payments made before the associated document has been created
type Prepayment struct {

	// See Prepayment Types
	Type string `json:"Type,omitempty"`

	// The date the prepayment is created YYYY-MM-DD
	Date string `json:"DateString,omitempty"`

	// See Contacts
	Contact Contact `json:"Contact"`

	// See Prepayment Status Codes
	Status string `json:"Status,omitempty"`

	// See Prepayment Line Amount Types
	LineAmountTypes string `json:"LineAmountTypes,omitempty"`

	// See Prepayment Line Items
	LineItems []LineItem `json:"LineItems,omitempty"`

	// The subtotal of the prepayment excluding taxes
	SubTotal float64 `json:"SubTotal,omitempty"`

	// The total tax on the prepayment
	TotalTax float64 `json:"TotalTax,omitempty"`

	// The total of the prepayment(subtotal + total tax)
	Total float64 `json:"Total,omitempty"`

	// UTC timestamp of last update to the prepayment
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty"`

	// Currency used for the prepayment
	CurrencyCode string `json:"CurrencyCode,omitempty"`

	// Xero generated unique identifier
	PrepaymentID string `json:"PrepaymentID,omitempty"`

	// The currency rate for a multicurrency prepayment. If no rate is specified, the XE.com day rate is used
	CurrencyRate float64 `json:"CurrencyRate,omitempty"`

	// The remaining credit balance on the prepayment
	RemainingCredit float64 `json:"RemainingCredit,omitempty"`

	// See Allocations
	Allocations []Allocation `json:"Allocations,omitempty"`

	// boolean to indicate if a prepayment has an attachment
	HasAttachments bool `json:"HasAttachments,omitempty"`
}

//Prepayments is a collection of Prepayments
type Prepayments struct {
	Prepayments []Prepayment `json:"Prepayments"`
}
