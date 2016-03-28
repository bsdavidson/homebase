package homebase

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"testing/iotest"
)

type fakeClient struct {
	StatusCode int
	Body       string
	Error      string
	BodyError  bool
}

func (f *fakeClient) Get(url string) (resp *http.Response, err error) {
	if f.Error != "" {
		return nil, errors.New(f.Error)
	}
	var reader io.Reader
	if f.BodyError {
		reader = iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBufferString("AB")))
	} else {
		reader = bytes.NewBufferString(f.Body)
	}
	body := ioutil.NopCloser(reader)
	resp = &http.Response{StatusCode: f.StatusCode, Body: body}
	return resp, nil
}

func TestGetPublicIP(t *testing.T) {
	client := &fakeClient{StatusCode: 200, Body: "  1.2.3.4  \n"}
	ip, err := GetPublicIP(client)
	if err != nil {
		t.Error("Expected error to be nil, was: ", err)
	}
	if ip.String() != "1.2.3.4" {
		t.Error("Expected IP to eq 1.2.3.4, got:", ip)
	}
}

func TestGetPublicIPBadStatus(t *testing.T) {
	client := &fakeClient{StatusCode: 400}
	ip, err := GetPublicIP(client)
	if err.Error() != "Bad status requesting IP: 400" {
		t.Error("Expected status code error, instead got: ", err)
	}
	if ip != nil {
		t.Error("Expected IP to be nil, instead got: ", ip)
	}
}

func TestGetPublicIPBadIP(t *testing.T) {
	client := &fakeClient{StatusCode: 200, Body: "I IZ BAD IP"}
	ip, err := GetPublicIP(client)
	if err.Error() != "Invalid IP: I IZ BAD IP" {
		t.Error("Expected bad IP error, instead got: ", err)
	}
	if ip != nil {
		t.Error("Expected IP to be nil, instead got: ", ip)
	}
}

func TestGetPublicIPGetError(t *testing.T) {
	client := &fakeClient{Error: "I HATE THE INTERNET"}
	ip, err := GetPublicIP(client)
	if err.Error() != "Error requesting IP: I HATE THE INTERNET" {
		t.Error("Expected GET error, instead got: ", err)
	}
	if ip != nil {
		t.Error("Expected IP to be nil, instead got: ", ip)
	}
}

func TestGetPublicIPGetBodyError(t *testing.T) {
	client := &fakeClient{StatusCode: 200, BodyError: true}
	ip, err := GetPublicIP(client)
	if err.Error() != "Error reading body: timeout" {
		t.Error("Expected read error, instead got: ", err)
	}
	if ip != nil {
		t.Error("Expected IP to be nil, instead got: ", ip)
	}
}
