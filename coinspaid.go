package coinspaid

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://app.coinspaid.com/api/v2/"
)

// Client manages communication with the Coinspaid API.
type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    *url.URL
	httpClient *http.Client
}

// ErrorResponse holds the error messages received from the API
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"error"`
	Code     string `json:"code"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v - %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Code)
}

// ValidationErrorResponse holds the error messages received from the API for validation errors
type ValidationErrorResponse struct {
	Response *http.Response
	Errors   map[string]string `json:"errors"`
}

func (r *ValidationErrorResponse) Error() string {
	return fmt.Sprintf("%v %v - %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Errors)
}

// NewClient returns a new instance of the Coinspaid client with the provided options
func NewClient(apiKey string, apiSecret string) *Client {

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	baseURL, _ := url.Parse(baseURL)

	return &Client{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (client *Client) doRequest(req *http.Request, v interface{}) (*http.Response, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	err = checkResponse(res)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&v)

	return res, err
}

// Address holds the data returned from the API
type Address struct {
	ID        int    `json:"id"`
	Currency  string `json:"currency"`
	ConvertTo string `json:"convert_to"`
	Address   string `json:"address"`
	Tag       string `json:"tag"`
	ForeignID string `json:"foreign_id"`
}

// UnmarshalJSON parses the request from server in the expected format
func (a *Address) UnmarshalJSON(data []byte) error {
	type Alias Address

	var temp struct {
		Data Alias `json:"data"`
	}

	err := json.Unmarshal(data, &temp)

	if err != nil {
		return nil
	}

	*a = Address(temp.Data)
	return nil
}

// TakeAddressInput specifies the parameters the TakeAddress method accepts.
type TakeAddressInput struct {
	// Your info for this address, will returned as reference in Address responses, example: user-id:2048
	ForeignID string `json:"foreign_id"`

	// ISO of currency to receive funds in, example: BTC
	Currency string `json:"currency"`
}

// TakeAddress Returns the address for depositing crypto
func (client *Client) TakeAddress(input *TakeAddressInput) (*Address, error) {

	relativeURL := &url.URL{Path: "addresses/take"}
	url := client.baseURL.ResolveReference(relativeURL)

	j, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(j))

	if err != nil {
		return nil, err
	}

	signedBody, err := client.createSignedRequestHeader(j)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Processing-Key", client.apiKey)
	req.Header.Set("X-Processing-Signature", signedBody)

	address := &Address{}

	_, err = client.doRequest(req, &address)

	if err != nil {
		return nil, err
	}

	return address, nil
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return errorResponse
	}

	if err == nil && len(body) > 0 {
		err := json.Unmarshal(body, errorResponse)
		if err != nil {
			errorResponse.Message = string(body)
		}
	}

	if r.StatusCode == http.StatusBadRequest {
		validationErrorResponse := &ValidationErrorResponse{Response: r}
		err = json.Unmarshal(body, validationErrorResponse)
		return validationErrorResponse
	}

	return errorResponse
}

func (client *Client) createSignedRequestHeader(body []byte) (response string, err error) {
	h := hmac.New(sha512.New, []byte(client.apiSecret))

	h.Write([]byte(body))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))
	return sha, nil
}
