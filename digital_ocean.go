package homebase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Record struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"data"`
}

type RecordsResponse struct {
	Records []Record `json:"domain_records"`
}

type HttpDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

var ErrorRecordNotFound = errors.New("Record not found")

func GetRecordByName(d HttpDoer, domainName, recordName, token string) (record *Record, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records", domainName), nil)
	if err != nil {
		return record, fmt.Errorf("Error creating request : %v", err)
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
		return nil, fmt.Errorf("Error parsing body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response requesting records: %d : %s", resp.StatusCode, string(jsonBytes))
	}

	recordsResponse := RecordsResponse{}
	if err := json.Unmarshal(jsonBytes, &recordsResponse); err != nil {
		return nil, fmt.Errorf("Error in json: %v", err)
	}

	for _, r := range recordsResponse.Records {
		if r.Name == recordName {
			return &r, nil
		}
	}

	return nil, ErrorRecordNotFound
}

func (r *Record) Save(client HttpDoer, domainName, token string) error {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("JSON error: %v", err)
	}
	url := fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records/%d", domainName, r.Id)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("Error creating request : %v", err)
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
		return fmt.Errorf("Error parsing body: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad response saving record: %d : %s", resp.StatusCode, string(jsonBytes))
	}
	return nil
}
