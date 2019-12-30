package accounting

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/helpers"
)

const (
	contactsURL = "https://api.xero.com/api.xro/2.0/Contacts"
)

//Contact is a debtor/customer or creditor/supplier in a Xero Organisation
type Contact struct {

	// Xero identifier
	ContactID string `json:"ContactID,omitempty" xml:"ContactID,omitempty"`

	// This can be updated via the API only i.e. This field is read only on the Xero contact screen, used to identify contacts in external systems (max length = 50). If the Contact Number is used, this is displayed as Contact Code in the Contacts UI in Xero.
	ContactNumber string `json:"ContactNumber,omitempty" xml:"ContactNumber,omitempty"`

	// A user defined account number. This can be updated via the API and the Xero UI (max length = 50)
	AccountNumber string `json:"AccountNumber,omitempty" xml:"AccountNumber,omitempty"`

	// Current status of a contact – see contact status types
	ContactStatus string `json:"ContactStatus,omitempty" xml:"ContactStatus,omitempty"`

	// Full name of contact/organisation (max length = 255)
	Name string `json:"Name,omitempty" xml:"Name,omitempty"`

	// First name of contact person (max length = 255)
	FirstName string `json:"FirstName,omitempty" xml:"FirstName,omitempty"`

	// Last name of contact person (max length = 255)
	LastName string `json:"LastName,omitempty" xml:"LastName,omitempty"`

	// Email address of contact person (umlauts not supported) (max length = 255)
	EmailAddress string `json:"EmailAddress,omitempty" xml:"EmailAddress,omitempty"`

	// Skype user name of contact
	SkypeUserName string `json:"SkypeUserName,omitempty" xml:"SkypeUserName,omitempty"`

	// See contact persons
	ContactPersons *[]ContactPerson `json:"ContactPersons,omitempty" xml:"ContactPersons>ContactPerson,omitempty"`

	// Bank account number of contact
	BankAccountDetails string `json:"BankAccountDetails,omitempty" xml:"BankAccountDetails,omitempty"`

	// Tax number of contact – this is also known as the ABN (Australia), GST Number (New Zealand), VAT Number (UK) or Tax ID Number (US and global) in the Xero UI depending on which regionalized version of Xero you are using (max length = 50)
	TaxNumber string `json:"TaxNumber,omitempty" xml:"TaxNumber,omitempty"`

	// Default tax type used for contact on AR Contacts
	AccountsReceivableTaxType string `json:"AccountsReceivableTaxType,omitempty" xml:"AccountsReceivableTaxType,omitempty"`

	// Default tax type used for contact on AP Contacts
	AccountsPayableTaxType string `json:"AccountsPayableTaxType,omitempty" xml:"AccountsPayableTaxType,omitempty"`

	// Store certain address types for a contact – see address types
	Addresses *[]Address `json:"Addresses,omitempty" xml:"Addresses>Address,omitempty"`

	// Store certain phone types for a contact – see phone types
	Phones *[]Phone `json:"Phones,omitempty" xml:"Phones>Phone,omitempty"`

	// true or false – Boolean that describes if a contact that has any AP Contacts entered against them. Cannot be set via PUT or POST – it is automatically set when an accounts payable Contact is generated against this contact.
	IsSupplier bool `json:"IsSupplier,omitempty" xml:"IsSupplier,omitempty"`

	// true or false – Boolean that describes if a contact has any AR Contacts entered against them. Cannot be set via PUT or POST – it is automatically set when an accounts receivable Contact is generated against this contact.
	IsCustomer bool `json:"IsCustomer,omitempty" xml:"IsCustomer,omitempty"`

	// Default currency for raising Contacts against contact
	DefaultCurrency string `json:"DefaultCurrency,omitempty" xml:"DefaultCurrency,omitempty"`

	// Store XeroNetworkKey for contacts.
	XeroNetworkKey string `json:"XeroNetworkKey,omitempty" xml:"XeroNetworkKey,omitempty"`

	// The default sales account code for contacts
	SalesDefaultAccountCode string `json:"SalesDefaultAccountCode,omitempty" xml:"SalesDefaultAccountCode,omitempty"`

	// The default purchases account code for contacts
	PurchasesDefaultAccountCode string `json:"PurchasesDefaultAccountCode,omitempty" xml:"PurchasesDefaultAccountCode,omitempty"`

	// The default sales tracking categories for contacts
	SalesTrackingCategories *[]TrackingCategory `json:"SalesTrackingCategories,omitempty" xml:"SalesTrackingCategories>SalesTrackingCategory,omitempty"`

	// The default purchases tracking categories for contacts
	PurchasesTrackingCategories *[]TrackingCategory `json:"PurchasesTrackingCategories,omitempty" xml:"PurchasesTrackingCategories>PurchaseTrackingCategory,omitempty"`

	// The name of the Tracking Category assigned to the contact under SalesTrackingCategories and PurchasesTrackingCategories
	TrackingCategoryName string `json:"TrackingCategoryName,omitempty" xml:"TrackingCategoryName,omitempty"`

	// The name of the Tracking Option assigned to the contact under SalesTrackingCategories and PurchasesTrackingCategories
	TrackingCategoryOption string `json:"TrackingCategoryOption,omitempty" xml:"TrackingCategoryOption,omitempty"`

	// UTC timestamp of last update to contact
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty" xml:"-"`

	// Displays which contact groups a contact is included in
	ContactGroups *[]ContactGroup `json:"ContactGroups,omitempty" xml:"ContactGroups>ContactGroup,omitempty"`

	// Website address for contact (read only)
	Website string `json:"Website,omitempty" xml:"-"`

	// batch payment details for contact (read only)
	BatchPayments BatchPayment `json:"BatchPayments,omitempty" xml:"-"`

	// The default discount rate for the contact (read only)
	Discount float64 `json:"Discount,omitempty" xml:"-"`

	// The raw AccountsReceivable(sales Contacts) and AccountsPayable(bills) outstanding and overdue amounts, not converted to base currency (read only)
	Balances Balances `json:"Balances,omitempty" xml:"-"`

	// A boolean to indicate if a contact has an attachment
	HasAttachments bool `json:"HasAttachments,omitempty" xml:"HasAttachments,omitempty"`
}

//Contacts contains a collection of Contacts
type Contacts struct {
	Contacts []Contact `json:"Contacts" xml:"Contact"`
}

//Balances are the raw AccountsReceivable(sales invoices) and AccountsPayable(bills)
//outstanding and overdue amounts, not converted to base currency
type Balances struct {
	AccountsReceivable Balance `json:"AccountsReceivable,omitempty" xml:"AccountsReceivable,omitempty"`
	AccountsPayable    Balance `json:"AccountsPayable,omitempty" xml:"AccountsPayable,omitempty"`
}

//Balance is the raw AccountsReceivable(sales invoices) and AccountsPayable(bills)
//outstanding and overdue amounts, not converted to base currency
type Balance struct {
	Outstanding float64 `json:"Outstanding,omitempty" xml:"Outstanding,omitempty"`
	Overdue     float64 `json:"Overdue,omitempty" xml:"Overdue,omitempty"`
}

//The Xero API returns Dates based on the .Net JSON date format available at the time of development
//We need to convert these to a more usable format - RFC3339 for consistency with what the API expects to recieve
func (c *Contacts) convertDates() error {
	var err error
	for n := len(c.Contacts) - 1; n >= 0; n-- {
		c.Contacts[n].UpdatedDateUTC, err = helpers.DotNetJSONTimeToRFC3339(c.Contacts[n].UpdatedDateUTC, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// unmarshalContact intermediate function used for apply the the changes in dates
// format
// TODO we can improve that overring the method Unmarshal
func unmarshalContact(contactResponseBytes []byte) (*Contacts, error) {
	var contactResponse *Contacts
	err := json.Unmarshal(contactResponseBytes, &contactResponse)
	if err != nil {
		return nil, err
	}

	err = contactResponse.convertDates()
	if err != nil {
		return nil, err
	}

	return contactResponse, err
}

// FindContacts will get all the contacts from Xero linked with the given
// tenantID
func FindContacts(cl *http.Client) (*Contacts, error) {
	request, err := http.NewRequest(http.MethodGet, contactsURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contactResponseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalContact(contactResponseBytes)
}

// FindContact will find the contact info with the given contactID
func FindContact(cl *http.Client, contactID uuid.UUID) (*Contact, error) {
	request, err := http.NewRequest(http.MethodGet, contactsURL+"/"+contactID.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contactResponseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	c, err := unmarshalContact(contactResponseBytes)
	if err != nil {
		return nil, err
	}
	if len(c.Contacts) > 0 {
		return &c.Contacts[0], nil
	}
	return nil, nil
}

// Create will create contacts with the given information
func (c *Contacts) Create(cl *http.Client) (*Contacts, error) {
	buf, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, contactsURL, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contactResponseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalContact(contactResponseBytes)
}

// Update will update the contact with the given criteria
func (c *Contact) Update(cl *http.Client) (*Contacts, error) {
	cn := Contacts{
		Contacts: []Contact{*c},
	}
	buf, err := json.Marshal(cn)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPut, contactsURL+"/"+c.ContactID, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contactResponseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalContact(contactResponseBytes)
}
