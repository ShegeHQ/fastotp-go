# FastOTP API Client for Go

This Go package provides a client for interacting with the FastOTP API.

## Installation

```bash
go get -u github.com/ShegeHQ/fastotp
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/ShegeHQ/fastotp"
)

func main() {
	// Initialize the client with your API key
	client := fastotp.Init("your_api_key")

	// Generate OTP
	generateReq := fastotp.GenerateOTPRequest{
		Type:        "alpha_numeric",
		Identifier:  "user123",
		Delivery:    map[string]string{"email": "user@example.com"},
		Validity:    5,
		TokenLength: 6,
	}

	generateResp, err := client.GenerateOTP(generateReq)
	if err != nil {
		log.Fatalf("Error generating OTP: %v", err)
	}

	fmt.Printf("Generated OTP: %+v\n", generateResp.OTP)

	// Validate OTP
	validateReq := fastotp.ValidateOTPRequest{
		Identifier: "user123",
		Token:      "123456",
	}

	validateResp, err := client.ValidateOTP(validateReq)
	if err != nil {
		log.Fatalf("Error validating OTP: %v", err)
	}

	fmt.Printf("Validated OTP: %+v\n", validateResp.OTP)

	// Get OTP details
	otpID := generateResp.OTP.ID
	otpDetails, err := client.GetOTP(otpID)
	if err != nil {
		log.Fatalf("Error fetching OTP details: %v", err)
	}

	fmt.Printf("OTP Details: %+v\n", otpDetails.OTP)
}
```

## Configuration

**ApiKey:** Your FastOTP API key.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or
submit a pull request.