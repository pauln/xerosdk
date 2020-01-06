package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/helpers"
)

const (
	itemURL = "https://api.xero.com/api.xro/2.0/Items"
)

//Item is something that is sold or purchased.  It can have inventory tracked or not tracked.
type Item struct {

	// User defined item code (max length = 30)
	Code string `json:"Code"`

	// The inventory asset account for the item. The account must be of type INVENTORY. The  COGSAccountCode in PurchaseDetails is also required to create a tracked item
	InventoryAssetAccountCode string `json:"InventoryAssetAccountCode"`

	// The name of the item (max length = 50)
	Name string `json:"Name,omitempty"`

	// Boolean value, defaults to true. When IsSold is true the item will be available on sales transactions in the Xero UI. If IsSold is updated to false then Description and SalesDetails values will be nulled.
	IsSold bool `json:"IsSold,omitempty"`

	// Boolean value, defaults to true. When IsPurchased is true the item is available for purchase transactions in the Xero UI. If IsPurchased is updated to false then PurchaseDescription and PurchaseDetails values will be nulled.
	IsPurchased bool `json:"IsPurchased,omitempty"`

	// The sales description of the item (max length = 4000)
	Description string `json:"Description,omitempty"`

	// The purchase description of the item (max length = 4000)
	PurchaseDescription string `json:"PurchaseDescription,omitempty"`

	// See Purchases & Sales
	PurchaseDetails PurchaseAndSaleDetails `json:"PurchaseDetails,omitempty"`

	// See Purchases & Sales
	SalesDetails PurchaseAndSaleDetails `json:"SalesDetails,omitempty"`

	// True for items that are tracked as inventory. An item will be tracked as inventory if the InventoryAssetAccountCode and COGSAccountCode are set.
	IsTrackedAsInventory bool `json:"IsTrackedAsInventory,omitempty"`

	// The value of the item on hand. Calculated using average cost accounting.
	TotalCostPool float64 `json:"TotalCostPool,omitempty"`

	// The quantity of the item on hand
	QuantityOnHand float64 `json:"QuantityOnHand,omitempty"`

	// Last modified date in UTC format
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty"`

	// The Xero identifier for an Item
	ItemID string `json:"ItemID,omitempty"`
}

//Items is a collection of Items
type Items struct {
	Items []Item `json:"Items"`
}

//PurchaseAndSaleDetails are Elements for Purchases and Sales
type PurchaseAndSaleDetails struct {
	//Unit Price of the item. By default UnitPrice is returned to two decimal places.  You can use 4 decimal places by adding the unitdp=4 querystring parameter to your request.
	UnitPrice float64 `json:"UnitPrice,omitempty"`

	//Default account code to be used for purchased/sale. Not applicable to the purchase details of tracked items
	AccountCode string `json:"AccountCode,omitempty"`

	//Cost of goods sold account. Only applicable to the purchase details of tracked items.
	COGSAccountCode string `json:"COGSAccountCode,omitempty"`

	//Used as an override if the default Tax Code for the selected AccountCode is not correct - see TaxTypes.
	TaxType string `json:"TaxType,omitempty"`
}

//The Xero API returns Dates based on the .Net JSON date format available at the time of development
//We need to convert these to a more usable format - RFC3339 for consistency with what the API expects to recieve
func (i *Items) convertDates() error {
	var err error
	for n := len(i.Items) - 1; n >= 0; n-- {
		i.Items[n].UpdatedDateUTC, err = helpers.DotNetJSONTimeToRFC3339(i.Items[n].UpdatedDateUTC, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalItem(itemResponseBytes []byte) (*Items, error) {
	var itemResponse *Items
	err := json.Unmarshal(itemResponseBytes, &itemResponse)
	if err != nil {
		return nil, err
	}

	err = itemResponse.convertDates()
	if err != nil {
		return nil, err
	}

	return itemResponse, err
}

// Create will create items given an Items struct
func (i *Items) Create(cl *http.Client) (*Items, error) {
	buf, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	itemsResponseBytes, err := helpers.Create(cl, itemURL, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalItem(itemsResponseBytes)
}

// Update will update an item given an Items struct
// This will only handle single item - you cannot update multiple items in a single call
func (i *Item) Update(cl *http.Client) (*Items, error) {
	its := Items{
		Items: []Item{*i},
	}
	buf, err := json.Marshal(its)
	if err != nil {
		return nil, err
	}
	itemsResponseBytes, err := helpers.Update(cl, itemURL+"/"+i.ItemID, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalItem(itemsResponseBytes)
}

// FindItems will get all items.
func FindItems(cl *http.Client, additionalHeaders map[string]string, queryParameters map[string]string) (*Items, error) {
	itemsResponseBytes, err := helpers.Find(cl, itemURL, additionalHeaders, queryParameters)
	if err != nil {
		return nil, err
	}

	return unmarshalItem(itemsResponseBytes)
}

//FindItem will get a single item - itemID must be a GUID for an item
func FindItem(cl *http.Client, itemID uuid.UUID) (*Item, error) {
	itemsResponseBytes, err := helpers.Find(cl, itemURL+"/"+itemID.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	items, err := unmarshalItem(itemsResponseBytes)
	if err != nil {
		return nil, err
	}
	return &items.Items[0], nil
}

//RemoveItem will get a single item - itemID must be a GUID for an item
func RemoveItem(cl *http.Client, itemID uuid.UUID) (*Items, error) {
	itemsResponseBytes, err := helpers.Remove(cl, itemURL+"/"+itemID.String())
	if err != nil {
		return nil, err
	}

	return unmarshalItem(itemsResponseBytes)
}
