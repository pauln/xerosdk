package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/quickaco/xerosdk/helpers"
)

const (
	contactGroupsURL = "https://api.xero.com/api.xro/2.0/ContactGroups"
)

//ContactGroup is a way of organising Contacts into groups
type ContactGroup struct {

	// The Name of the contact group. Required when creating a new contact group
	Name string `json:"Name,omitempty"`

	// The Status of a contact group. To delete a contact group update the status to DELETED. Only contact groups with a status of ACTIVE are returned on GETs.
	Status string `json:"Status,omitempty"`

	// The Xero identifier for an contact group – specified as a string following the endpoint name. e.g. /297c2dc5-cc47-4afd-8ec8-74990b8761e9
	ContactGroupID string `json:"ContactGroupID,omitempty"`

	// The ContactID and Name of Contacts in a contact group. Returned on GETs when the ContactGroupID is supplied in the URL.
	Contacts []Contact `json:"Contacts,omitempty"`
}

//ContactGroups is a collection of ContactGroups
type ContactGroups struct {
	ContactGroups []ContactGroup `json:"ContactGroups"`
}

func unmarshalContactGroup(contactGroupResponseBytes []byte) (*ContactGroups, error) {
	var contactGroupResponse *ContactGroups
	err := json.Unmarshal(contactGroupResponseBytes, &contactGroupResponse)
	if err != nil {
		return nil, err
	}

	return contactGroupResponse, err
}

// FindContactGroups will get all contactGroups
func FindContactGroups(cl *http.Client) (*ContactGroups, error) {
	contactGroupsBytes, err := helpers.Find(cl, contactGroupsURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return unmarshalContactGroup(contactGroupsBytes)
}

// FindContactGroup will get a single contactGroup - contactGroupID must be a GUID for an contactGroup
func FindContactGroup(cl *http.Client, contactGroupID uuid.UUID) (*ContactGroups, error) {
	contactGroupsBytes, err := helpers.Find(cl, contactGroupsURL+"/"+contactGroupID.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	return unmarshalContactGroup(contactGroupsBytes)
}

// RemoveContactGroup will get a single contactGroup - contactGroupID must be a GUID for an contactGroup
func RemoveContactGroup(cl *http.Client, contactGroupID uuid.UUID) (*ContactGroups, error) {
	contactGroupsBytes, err := helpers.Remove(cl, contactGroupsURL+"/"+contactGroupID.String())
	if err != nil {
		return nil, err
	}

	return unmarshalContactGroup(contactGroupsBytes)
}

//Create will create contactGroups given an ContactGroups struct
func (c *ContactGroups) Create(cl *http.Client) (*ContactGroups, error) {
	buf, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	contactGroupBytes, err := helpers.Create(cl, contactGroupsURL, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalContactGroup(contactGroupBytes)
}

//Update will update an contactGroup given an ContactGroups struct
//This will only handle single contactGroup - you cannot update multiple contactGroups in a single call
func (c *ContactGroup) Update(cl *http.Client) (*ContactGroups, error) {
	cg := ContactGroups{
		ContactGroups: []ContactGroup{*c},
	}
	buf, err := json.Marshal(cg)
	if err != nil {
		return nil, err
	}
	contactGroupBytes, err := helpers.Update(cl, contactGroupsURL+"/"+c.ContactGroupID, buf)
	if err != nil {
		return nil, err
	}

	return unmarshalContactGroup(contactGroupBytes)
}
