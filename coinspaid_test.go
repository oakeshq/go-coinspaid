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
)

func TestClient(t *testing.T) {
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
