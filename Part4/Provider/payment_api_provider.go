package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/m/httpclient"
)

type PaymentAPIProvider interface {
	CallPaymentAPI(request PaymentRequest) (*PaymentResponse, *SystemError)
}

type PaymentAPI struct {
	httpClient httpclient.HTTPClient
}

func NewPaymentProvider(client httpclient.HTTPClient) PaymentAPIProvider {
	return &PaymentAPI{httpClient: client}
}

type PaymentRequest struct {
	CardNo          string
	CardHolderName  string
	Amount          string
	RequestNumber   string
	RequestDateTime string
}

type PaymentResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type SystemError struct {
	ErrorCode    string
	ErrorMessage string
}

func (p PaymentAPI) CallPaymentAPI(request PaymentRequest) (*PaymentResponse, *SystemError) {

	req, err := http.NewRequest("POST", "https://external.payment.api", nil)
	if err != nil {
		return nil, &SystemError{ErrorCode: "ERR1", ErrorMessage: "Failed to create request"}
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, &SystemError{ErrorCode: "ERR2", ErrorMessage: "HTTP request failed"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &SystemError{ErrorCode: "ERR3", ErrorMessage: fmt.Sprintf("external api returned status %d", resp.StatusCode)}
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, &SystemError{ErrorCode: "ERR4", ErrorMessage: "Failed to decode response"}
	}

	return &paymentResp, nil
}
