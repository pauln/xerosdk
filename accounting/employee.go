package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/quickaco/xerosdk/helpers"
)

const (
	employeeURL = "https://api.xero.com/api.xro/2.0/Employees"
)

//Employee is for the deprecated Pay run feature.
type Employee struct {

	// The Xero identifier for an employee e.g. 297c2dc5-cc47-4afd-8ec8-74990b8761e9
	EmployeeID string `json:"EmployeeID,omitempty"`

	// Current status of an employee – see contact status types
	Status string `json:"Status,omitempty"`

	// First name of an employee (max length = 255)
	FirstName string `json:"FirstName,omitempty"`

	// Last name of an employee (max length = 255)
	LastName string `json:"LastName,omitempty"`

	// ExternalLink of an employee
	ExternalLink ExternalLink `json:"ExternalLink,omitempty"`
}

// Employees is for encapsulate the response from Employee
type Employees struct {
	Employess []Employee `json:"Employees,omitempty"`
}

// FindEmployees will find the info about employees
func FindEmployees(cl *http.Client, queryParameters map[string]string) (em *Employees, err error) {
	employeeResponseBytes, err := helpers.Find(cl, employeeURL, nil, queryParameters)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(employeeResponseBytes, &em); err != nil {
		return nil, err
	}
	return em, nil
}

// Create will create employees with the given information
func (e *Employees) Create(cl *http.Client) (em *Employees, err error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	employeeResponseBytes, err := helpers.Create(cl, employeeURL, buf)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(employeeResponseBytes, &em); err != nil {
		return nil, err
	}
	return em, nil
}

// Update will update employees with the given criteria
func (e *Employee) Update(cl *http.Client) (em *Employees, err error) {
	es := Employees{
		Employess: []Employee{*e},
	}
	buf, err := json.Marshal(es)
	if err != nil {
		return nil, err
	}
	employeeResponseBytes, err := helpers.Update(cl, employeeURL+"/"+e.EmployeeID, buf)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(employeeResponseBytes, &em); err != nil {
		return nil, err
	}
	return em, nil
}
