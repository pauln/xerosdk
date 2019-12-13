package connection

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
)

const (
	connectionsURL = "https://api.xero.com/connections"
)

// Tenant type will keep information about the Xero tenant
type Tenant struct {
	TenantID   uuid.UUID `json:"tenantId,omitempty"`
	TenantType string    `json:"tenantType,omitempty"`
}

// GetTenants will return the value of the getting information from xero
func GetTenants(c *http.Client) ([]Tenant, error) {
	tenants := []Tenant{}
	request, err := http.NewRequest(http.MethodGet, connectionsURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")

	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&tenants); err != nil {
		return nil, err
	}
	return tenants, nil
}
