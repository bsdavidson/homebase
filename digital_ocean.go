package homebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Record is a Digital Ocean domain record.
type Record struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"data"`
}

type recordsResponse struct {
	Records []Record `json:"domain_records"`
}

// HTTPDoer is an interface that implements http.Client's Do method so that mock
// objects can be passed in for tests.
type HTTPDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// ErrRecordNotFound is returned when a matching record cannot be found.
var ErrRecordNotFound = errors.New("Record not found")

// GetRecordByName returns a record of a subdomain for a given Digital Ocean
// domain. If a record cannot be found with the passed name, ErrRecordNotFound
// will be returned instead.
func GetRecordByName(d HTTPDoer, domainName, recordName, token string) (record *Record, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records", domainName), nil)
	if err != nil {
		return record, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := d.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making the request: %v", err)
	}
	defer resp.Body.Close()
	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response requesting records: %d: %s", resp.StatusCode, string(jsonBytes))
	}

	recordsResponse := recordsResponse{}
	if err := json.Unmarshal(jsonBytes, &recordsResponse); err != nil {
		return nil, fmt.Errorf("Error in json: %v", err)
	}

	for _, r := range recordsResponse.Records {
		if r.Name == recordName {
			return &r, nil
		}
	}

	return nil, ErrRecordNotFound
}

// Save a record to Digital Ocean.
func (r *Record) Save(client HTTPDoer, domainName, token string) error {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("JSON error: %v", err)
	}
	url := fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records/%d", domainName, r.ID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making the request: %v", err)
	}
	defer resp.Body.Close()
	jsonBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading body: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad response saving record: %d: %s", resp.StatusCode, string(jsonBytes))
	}
	return nil
}
