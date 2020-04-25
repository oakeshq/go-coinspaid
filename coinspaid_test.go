package coinspaid

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	okResponse = `{
		"data": {
			"id": 1,
			"currency": "EUR",
			"convert_to": "EUR",
			"address": "12983h13ro1hrt24it432t",
			"tag": "tag-123",
			"foreign_id": "user-id:2048"
		}
	}`

	invalidAuthResponse = `{
		"error": "Bad key header",
		"code": "bad_header_key"
	}`

	badRequestResponse = `{
		"errors": {
			"foreign_id": "The foreign id field is required."
		}
	}`
)

func TestTakeAddress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(okResponse))
	}))

	defer server.Close()

	baseURL, _ := url.Parse(server.URL)

	api := Client{
		apiKey:     "key",
		apiSecret:  "secret",
		httpClient: server.Client(),
		baseURL:    baseURL,
	}

	takeAddressInput := &TakeAddressInput{
		ForeignID: "user-id:2048",
		Currency:  "EUR",
	}

	address, err := api.TakeAddress(takeAddressInput)

	assert.Nil(t, err)
	assert.Equal(t, takeAddressInput.Currency, address.Currency)
}

func TestClientWithInvalidAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(invalidAuthResponse))
	}))

	defer server.Close()

	baseURL, _ := url.Parse(server.URL)

	api := Client{
		apiKey:     "invalid",
		apiSecret:  "invalid",
		httpClient: server.Client(),
		baseURL:    baseURL,
	}

	takeAddressInput := &TakeAddressInput{
		ForeignID: "user-id:2048",
		Currency:  "EUR",
	}

	_, err := api.TakeAddress(takeAddressInput)

	assert.NotNil(t, err)
	assert.Equal(t, "bad_header_key", err.(*ErrorResponse).Code)
}

func TestClientWithBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(badRequestResponse))
	}))

	defer server.Close()

	baseURL, _ := url.Parse(server.URL)

	api := Client{
		apiKey:     "invalid",
		apiSecret:  "invalid",
		httpClient: server.Client(),
		baseURL:    baseURL,
	}

	takeAddressInput := &TakeAddressInput{
		Currency: "INEXISTENT",
	}

	_, err := api.TakeAddress(takeAddressInput)

	assert.NotNil(t, err)
	assert.NotNil(t, err.(*ValidationErrorResponse).Errors)
}
