package homebase

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
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
		// When body error is specified, we want to simulate an error being returned while reading the body.
		// TimeoutReader will error on the second read.
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
	if assert.NoError(t, err) {
		assert.Equal(t, "1.2.3.4", ip.String())
	}
}

func TestGetPublicIPBadStatus(t *testing.T) {
	client := &fakeClient{StatusCode: 400}
	ip, err := GetPublicIP(client)
	if assert.Error(t, err) {
		assert.Equal(t, "Bad status requesting IP: 400", err.Error())
	}
	assert.Nil(t, ip)
}

func TestGetPublicIPBadIP(t *testing.T) {
	client := &fakeClient{StatusCode: 200, Body: "I IZ BAD IP"}
	ip, err := GetPublicIP(client)
	if assert.Error(t, err) {
		assert.Equal(t, "Invalid IP: I IZ BAD IP", err.Error())
	}
	assert.Nil(t, ip)
}

func TestGetPublicIPGetError(t *testing.T) {
	client := &fakeClient{Error: "I HATE THE INTERNET"}
	ip, err := GetPublicIP(client)
	if assert.Error(t, err) {
		assert.Equal(t, "Error requesting IP: I HATE THE INTERNET", err.Error())
	}
	assert.Nil(t, ip)
}

func TestGetPublicIPGetBodyError(t *testing.T) {
	client := &fakeClient{StatusCode: 200, BodyError: true}
	ip, err := GetPublicIP(client)
	if assert.Error(t, err) {
		assert.Equal(t, "Error reading body: timeout", err.Error())
	}
	assert.Nil(t, ip)
}
