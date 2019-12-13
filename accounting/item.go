package accounting

//Item is something that is sold or purchased.  It can have inventory tracked or not tracked.
type Item struct {

	// User defined item code (max length = 30)
	Code string `json:"Code" xml:"Code"`

	// The inventory asset account for the item. The account must be of type INVENTORY. The  COGSAccountCode in PurchaseDetails is also required to create a tracked item
	InventoryAssetAccountCode string `json:"InventoryAssetAccountCode" xml:"InventoryAssetAccountCode"`

	// The name of the item (max length = 50)
	Name string `json:"Name,omitempty" xml:"Name,omitempty"`

	// Boolean value, defaults to true. When IsSold is true the item will be available on sales transactions in the Xero UI. If IsSold is updated to false then Description and SalesDetails values will be nulled.
	IsSold bool `json:"IsSold,omitempty" xml:"IsSold,omitempty"`

	// Boolean value, defaults to true. When IsPurchased is true the item is available for purchase transactions in the Xero UI. If IsPurchased is updated to false then PurchaseDescription and PurchaseDetails values will be nulled.
	IsPurchased bool `json:"IsPurchased,omitempty" xml:"IsPurchased,omitempty"`

	// The sales description of the item (max length = 4000)
	Description string `json:"Description,omitempty" xml:"Description,omitempty"`

	// The purchase description of the item (max length = 4000)
	PurchaseDescription string `json:"PurchaseDescription,omitempty" xml:"PurchaseDescription,omitempty"`

	// See Purchases & Sales
	PurchaseDetails PurchaseAndSaleDetails `json:"PurchaseDetails,omitempty" xml:"PurchaseDetails,omitempty"`

	// See Purchases & Sales
	SalesDetails PurchaseAndSaleDetails `json:"SalesDetails,omitempty" xml:"SalesDetails,omitempty"`

	// True for items that are tracked as inventory. An item will be tracked as inventory if the InventoryAssetAccountCode and COGSAccountCode are set.
	IsTrackedAsInventory bool `json:"IsTrackedAsInventory,omitempty" xml:"-"`

	// The value of the item on hand. Calculated using average cost accounting.
	TotalCostPool float64 `json:"TotalCostPool,omitempty" xml:"-"`

	// The quantity of the item on hand
	QuantityOnHand float64 `json:"QuantityOnHand,omitempty" xml:"-"`

	// Last modified date in UTC format
	UpdatedDateUTC string `json:"UpdatedDateUTC,omitempty" xml:"-"`

	// The Xero identifier for an Item
	ItemID string `json:"ItemID,omitempty" xml:"ItemID,omitempty"`
}

//Items is a collection of Items
type Items struct {
	Items []Item `json:"Items" xml:"Item"`
}

//PurchaseAndSaleDetails are Elements for Purchases and Sales
type PurchaseAndSaleDetails struct {
	//Unit Price of the item. By default UnitPrice is returned to two decimal places.  You can use 4 decimal places by adding the unitdp=4 querystring parameter to your request.
	UnitPrice float64 `json:"UnitPrice,omitempty" xml:"UnitPrice,omitempty"`

	//Default account code to be used for purchased/sale. Not applicable to the purchase details of tracked items
	AccountCode string `json:"AccountCode,omitempty" xml:"AccountCode,omitempty"`

	//Cost of goods sold account. Only applicable to the purchase details of tracked items.
	COGSAccountCode string `json:"COGSAccountCode,omitempty" xml:"COGSAccountCode,omitempty"`

	//Used as an override if the default Tax Code for the selected AccountCode is not correct - see TaxTypes.
	TaxType string `json:"TaxType,omitempty" xml:"TaxType,omitempty"`
}
