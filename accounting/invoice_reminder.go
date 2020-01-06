package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/quickaco/xerosdk/helpers"
)

const (
	invoiceRemindersURL = "https://api.xero.com/api.xro/2.0/InvoiceReminders/Settings"
)

// InvoiceReminder will keep information about invoicing settings
type InvoiceReminder struct {
	Enabled bool `json:"Enabled,omitempty"`
}

// InvoiceReminders will keep information about the response from Xero
type InvoiceReminders struct {
	InvoiceReminders []InvoiceReminder `json:"InvoiceReminders,omitempty"`
}

// FindInvoiceReminders will get all the invoice reminders from Xero
func FindInvoiceReminders(cl *http.Client) (ir *InvoiceReminders, err error) {
	invoiceRemindersBytes, err := helpers.Find(cl, invoiceRemindersURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(invoiceRemindersBytes, &ir); err != nil {
		return nil, err
	}
	return ir, nil
}
