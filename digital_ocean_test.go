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

type fakeDoer struct {
	statusCode int
	body       string
	err        string
	bodyError  bool
}

func (f *fakeDoer) Do(req *http.Request) (resp *http.Response, err error) {
	if f.err != "" {
		return nil, errors.New(f.err)
	}
	var reader io.Reader
	if f.bodyError {
		// When body error is specified, we want to simulate an error being returned while reading the body.
		// TimeoutReader will error on the second read.
		reader = iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBufferString("AB")))
	} else {
		reader = bytes.NewBufferString(f.body)
	}
	body := ioutil.NopCloser(reader)
	resp = &http.Response{StatusCode: f.statusCode, Body: body}
	return resp, nil
}

func TestGetRecordByName(t *testing.T) {
	c := &fakeDoer{statusCode: 200, body: `{
    "domain_records": [
      {
        "id": 1,
        "type": "A",
        "name": "subdomain",
        "data": "1.2.3.4"
      }
    ]
  }`}

	record, err := GetRecordByName(c, "example.com", "subdomain", "12345")
	if assert.NoError(t, err) {
		assert.NotNil(t, record)
		assert.Equal(t, 1, record.ID)
		assert.Equal(t, "A", record.Type)
		assert.Equal(t, "subdomain", record.Name)
		assert.Equal(t, "1.2.3.4", record.Data)
	}

	record, err = GetRecordByName(c, "example.com", "notfound", "12345")
	if assert.Error(t, err) {
		assert.Equal(t, ErrRecordNotFound, err)
		assert.Nil(t, record)
	}

}

func TestGetRecordByNameErrors(t *testing.T) {
	type testFixture struct {
		fakeDoer
		errorString string
	}
	testFixtures := []testFixture{
		testFixture{
			errorString: "Error in json: unexpected end of JSON input",
			fakeDoer:    fakeDoer{statusCode: 200, body: ""},
		},
		testFixture{
			errorString: "Bad response requesting records: 404: Not Found",
			fakeDoer:    fakeDoer{statusCode: 404, body: "Not Found"},
		},
		testFixture{
			errorString: "Error reading body: timeout",
			fakeDoer:    fakeDoer{statusCode: 200, bodyError: true},
		},
		testFixture{
			errorString: "Error making the request: Not awesome error",
			fakeDoer:    fakeDoer{err: "Not awesome error"},
		},
	}

	for _, tf := range testFixtures {
		record, err := GetRecordByName(&tf.fakeDoer, "example.com", "subdomain", "12345")
		if assert.Error(t, err) {
			assert.Equal(t, tf.errorString, err.Error())
			assert.Nil(t, record)
		}
	}
}

func TestRecordSave(t *testing.T) {
	record := &Record{
		ID:   1,
		Type: "A",
		Name: "subdomain",
		Data: "1.2.3.4",
	}

	fc := &fakeDoer{statusCode: 200}

	err := record.Save(fc, "example.com", "12345")
	assert.NoError(t, err)
}

func TestRecordSaveErrors(t *testing.T) {
	type testFixture struct {
		fakeDoer
		errorString string
	}
	testFixtures := []testFixture{
		testFixture{
			errorString: "Bad response saving record: 404: Not Found",
			fakeDoer:    fakeDoer{statusCode: 404, body: "Not Found"},
		},
		testFixture{
			errorString: "Error reading body: timeout",
			fakeDoer:    fakeDoer{statusCode: 200, bodyError: true},
		},
		testFixture{
			errorString: "Error making the request: Not awesome error",
			fakeDoer:    fakeDoer{err: "Not awesome error"},
		},
	}

	record := &Record{
		ID:   1,
		Type: "A",
		Name: "subdomain",
		Data: "1.2.3.4",
	}

	for _, tf := range testFixtures {
		err := record.Save(&tf.fakeDoer, "example.com", "12345")
		if assert.Error(t, err) {
			assert.Equal(t, tf.errorString, err.Error())
		}
	}
}
