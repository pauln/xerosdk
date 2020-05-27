package helpers

import "encoding/json"

// Error is a type that tries to decode the Xero API error, couldn't find anything
// on the documentation that give me a clear vision about how Xero is managing
// the errors. https://developer.xero.com/documentation/api/http-response-codes
// Added this structure from a real api call
type Error struct {
	Title    string
	Status   int
	Detail   string
	Instance string
}

// DecodeError will try to match the given error into the Error type, if it
// fails will return a generic Error object
func DecodeError(buf []byte) Error {
	var e Error
	if err := json.Unmarshal(buf, &e); err != nil {
		return Error{
			Title:  "Unknown error",
			Status: 500,
			Detail: "Error decoding the xero error response",
		}
	}
	return e
}
