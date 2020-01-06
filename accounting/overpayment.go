package accounting

//Overpayment is used when a debtor overpays an invoice
type Overpayment struct {

	// See Overpayment Types
	Type string `json:"Type,omitempty"`

	// The date the overpayment is created YYYY-MM-DD
	Date string `json:"DateString,omitempty"`

	// See Contacts
	Contact Contact `json:"Contact"`

	// See Overpayment Status Codes
	Status string `json:"Status,omitempty"`

	// See Overpayment Line Amount Types
	LineAmountTypes string `json:"LineAmountTypes,omitempty"`

	// See Overpayment Line Items
	LineItems []LineItem `json:"LineItems,omitempty"`

	// The subtotal of the overpayment excluding taxes
	SubTotal float64 `json:"SubTotal,omitempty"`

	// The total tax on the overpayment
	TotalTax float64 `json:"TotalTax,omitempty"`

	// The total of the overpayment (subtotal + total tax)
	Total float64 `json:"Total,omitempty"`

	// UTC timestamp of last update to the overpayment
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty"`

	// Currency used for the overpayment
	CurrencyCode string `json:"CurrencyCode,omitempty"`

	// Xero generated unique identifier
	OverpaymentID string `json:"OverpaymentID,omitempty"`

	// The currency rate for a multicurrency overpayment. If no rate is specified, the XE.com day rate is used
	CurrencyRate float64 `json:"CurrencyRate,omitempty"`

	// The remaining credit balance on the overpayment
	RemainingCredit float64 `json:"RemainingCredit,omitempty"`

	// See Allocations
	Allocations []Allocation `json:"Allocations,omitempty"`

	// See Payments
	Payments []Payment `json:"Payments,omitempty"`

	// boolean to indicate if a overpayment has an attachment
	HasAttachments bool `json:"HasAttachments,omitempty"`
}

//Overpayments is a collection of Overpayments
type Overpayments struct {
	Overpayments []Overpayment `json:"Overpayments"`
}
