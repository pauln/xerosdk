package helpers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

// Find function encapsulate all the GET method calls to Xero API
func Find(cl *http.Client, endpoint string, additionalHeaders map[string]string, queryParameters map[string]string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	values := request.URL.Query()
	for key, value := range queryParameters {
		values.Add(key, value)
	}
	request.URL.RawQuery = values.Encode()
	for key, value := range additionalHeaders {
		request.Header.Add(key, value)
	}

	return process(cl, request)
}

// Create function encapsulate all the POST method calls to Xero API
func Create(cl *http.Client, endpoint string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	return process(cl, request)
}

// Update function encapsulate all the PUT method calls to Xero API
func Update(cl *http.Client, endpoint string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	return process(cl, request)
}

// Remove function encapsulate all the DELETE method calls to Xero API
func Remove(cl *http.Client, endpoint string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return process(cl, request)
}

func process(cl *http.Client, request *http.Request) ([]byte, error) {
	request.Header.Add("Accept", "application/json")
	response, err := cl.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(string(responseBytes))
	}
	return responseBytes, nil
}
