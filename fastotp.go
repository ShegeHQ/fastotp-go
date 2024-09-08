package fastotp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func Init(apiKey string) *Client {
	return &Client{
		BaseURL: "https://api.fastotp.co",
		APIKey:  apiKey,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type OTP struct {
	ID              string                 `json:"id"`
	Identifier      string                 `json:"identifier"`
	Type            string                 `json:"type"`
	Status          string                 `json:"status"`
	DeliveryMethods []string               `json:"delivery_methods"`
	DeliveryDetails map[string]interface{} `json:"delivery_details"`
	ExpiresAt       string                 `json:"expires_at"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

type GenerateOTPRequest struct {
	Type        string            `json:"type"`
	Identifier  string            `json:"identifier"`
	Delivery    map[string]string `json:"delivery"`
	Validity    int               `json:"validity"`
	TokenLength int               `json:"token_length"`
}

type GenerateOTPResponse struct {
	OTP OTP `json:"otp"`
}

type ValidateOTPRequest struct {
	Identifier string `json:"identifier"`
	Token      string `json:"token"`
}

type ValidateOTPResponse struct {
	OTP OTP `json:"otp"`
}

type ErrorResponse struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func (c *Client) GenerateOTP(req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	url := fmt.Sprintf("%s/generate", c.BaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", c.APIKey)

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(response)
	}

	var otpResponse GenerateOTPResponse
	if err := json.NewDecoder(response.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	return &otpResponse, nil
}

func (c *Client) ValidateOTP(req ValidateOTPRequest) (*ValidateOTPResponse, error) {
	url := fmt.Sprintf("%s/validate", c.BaseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", c.APIKey)

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(response)
	}

	var otpResponse ValidateOTPResponse
	if err := json.NewDecoder(response.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	return &otpResponse, nil
}

func (c *Client) GetOTP(id string) (*GenerateOTPResponse, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, id)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("x-api-key", c.APIKey)

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(response)
	}

	var otpResponse GenerateOTPResponse
	if err := json.NewDecoder(response.Body).Decode(&otpResponse); err != nil {
		return nil, err
	}

	return &otpResponse, nil
}

func parseErrorResponse(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}

	return fmt.Errorf("error: %s, details: %v", errResp.Message, errResp.Errors)
}
