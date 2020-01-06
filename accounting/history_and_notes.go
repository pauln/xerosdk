package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/quickaco/xerosdk/helpers"
)

const (
	historyRecordURL = "https://api.xero.com/api.xro/2.0/"
)

// HistoryRecord is a record of monies transferred from one bank account to another
type HistoryRecord struct {

	// The type of change recorded against the document
	Changes string `json:"Changes,omitempty"`

	// UTC date that the history record was created
	DateUTC string `json:"DateUTC,omitempty"`

	// The user responsible for the change ("System Generated" when the change happens via API)
	User string `json:"User,omitempty"`

	// The Bank Transaction ID for the source account
	Details string `json:"Details"`
}

// HistoryRecords contains a collection of BankTransfers
type HistoryRecords struct {
	HistoryRecords []HistoryRecord `json:"HistoryRecords"`
}

//The Xero API returns Dates based on the .Net JSON date format available at the time of development
//We need to convert these to a more usable format - RFC3339 for consistency with what the API expects to recieve
func (h *HistoryRecords) convertDates() error {
	var err error
	for n := len(h.HistoryRecords) - 1; n >= 0; n-- {
		h.HistoryRecords[n].DateUTC, err = helpers.DotNetJSONTimeToRFC3339(h.HistoryRecords[n].DateUTC, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func unmarshalHistoryRecord(HistoryRecordResponseBytes []byte) (*HistoryRecords, error) {
	var historyRecordResponse *HistoryRecords
	err := json.Unmarshal(HistoryRecordResponseBytes, &historyRecordResponse)
	if err != nil {
		return nil, err
	}

	err = historyRecordResponse.convertDates()
	if err != nil {
		return nil, err
	}

	return historyRecordResponse, err
}

//FindHistoryAndNotes gets all history items and notes for a given type and ID.
//it is not supported on all endpoints.  See https://developer.xero.com/documentation/api/history-and-notes#SupportedDocs
func FindHistoryAndNotes(cl *http.Client, docType string, id string) (*HistoryRecords, error) {
	historyAndNotesBytes, err := helpers.Find(cl, historyRecordURL+docType+"/"+id+"/history", nil, nil)
	if err != nil {
		return nil, err
	}

	return unmarshalHistoryRecord(historyAndNotesBytes)
}

// Create will create History Records given a HistoryRecords struct and a docType and id
func (h *HistoryRecords) Create(cl *http.Client, docType string, id string) (*HistoryRecords, error) {
	buf, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	historyAndNotesBytes, err := helpers.Create(cl, historyRecordURL+docType+"/"+id+"/history", buf)
	if err != nil {
		return nil, err
	}

	return unmarshalHistoryRecord(historyAndNotesBytes)
}
