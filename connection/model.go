package connection

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/helpers"
)

const (
	connectionsURL = "https://api.xero.com/connections"
)

// Tenant type will keep information about the Xero tenant
type Tenant struct {
	ID         uuid.UUID `json:"id,omitempty"`
	TenantID   uuid.UUID `json:"tenantId,omitempty"`
	TenantType string    `json:"tenantType,omitempty"`
}

// GetTenants will return the value of the getting information from xero
func GetTenants(cl *http.Client) (tenants []Tenant, err error) {
	tenantResponseBytes, err := helpers.Find(cl, connectionsURL, nil, nil)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(tenantResponseBytes, &tenants); err != nil {
		return nil, err
	}
	return tenants, nil
}

// DeleteTenant will remove the connection with the given connectionID
func DeleteTenant(cl *http.Client, connectionID uuid.UUID) error {
	_, err := helpers.Remove(cl, connectionsURL+"/"+connectionID.String())
	if err != nil {
		return err
	}
	return nil
}
