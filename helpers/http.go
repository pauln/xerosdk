package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Find function encapsulate all the GET method calls to Xero API
func Find(cl *http.Client, endpoint string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

// Create function encapsulate all the POST method calls to Xero API
func Create(cl *http.Client, endpoint string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

// Update function encapsulate all the PUT method calls to Xero API
func Update(cl *http.Client, endpoint string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
