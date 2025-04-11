package cashrampsdkgo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rockets-hq/cashramp-sdk-go/queries"
	"github.com/rockets-hq/cashramp-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func mockGraphQLServer(t *testing.T, mockResponse []byte, expectedStatusCode int, expectedAuthHeader bool, expectedBodyContains ...string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		if expectedAuthHeader {
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		}

		if len(expectedBodyContains) > 0 {
			buf := new(strings.Builder)
			_, err := io.Copy(buf, r.Body)
			assert.NoError(t, err)
			bodyStr := buf.String()
			for _, substr := range expectedBodyContains {
				assert.Contains(t, bodyStr, substr)
			}
		}

		w.WriteHeader(expectedStatusCode)
		w.Write(mockResponse)
	}))
}

func dummyClient(server *httptest.Server) *Client {
	return &Client{
		apiUrl:     server.URL,
		secretKey:  "dummy-secret",
		httpClient: server.Client(),
	}
}

func createMockGraphQLResponse(t *testing.T, operationName string, result any, graphqlErrors ...string) []byte {
	response := map[string]any{"data": map[string]any{}}
	if result != nil {
		response["data"].(map[string]any)[operationName] = result
	}

	if len(graphqlErrors) > 0 {
		errs := []map[string]string{}
		for _, msg := range graphqlErrors {
			errs = append(errs, map[string]string{"message": msg})
		}
		response["errors"] = errs
		if result == nil {
			response["data"] = nil
		}
	}

	responseBytes, err := json.Marshal(response)
	assert.NoError(t, err)
	return responseBytes
}

func TestInitialiseClient(t *testing.T) {
	t.Run("Valid test environment", func(t *testing.T) {
		client, err := InitialiseClient("test", "dummy-secret")
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "test", client.env)
		assert.Contains(t, client.apiUrl, "staging.api.useaccrue.com")
		assert.Equal(t, "dummy-secret", client.secretKey)
	})

	t.Run("Valid live environment", func(t *testing.T) {
		client, err := InitialiseClient("live", "dummy-secret")
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "live", client.env)
		assert.Contains(t, client.apiUrl, "api.useaccrue.com")
		assert.NotContains(t, client.apiUrl, "staging")
		assert.Equal(t, "dummy-secret", client.secretKey)
	})

	t.Run("Invalid environment", func(t *testing.T) {
		_, err := InitialiseClient("invalid-env", "dummy-secret")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid env")
	})

	t.Run("Missing secret key", func(t *testing.T) {
		_, err := InitialiseClient("test", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Please provide your API secret key")
	})
}

func TestSendRequestSuccess(t *testing.T) {
	operationName := "account"
	mockResult := map[string]string{
		"accountBalance": "10",
		"depositAddress": "0x123",
	}
	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName)
	defer server.Close()

	client := dummyClient(server)

	resp, err := client.SendRequest(operationName, queries.ACCOUNT, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Empty(t, resp.Error)
	assert.NotNil(t, resp.Result)

	resultMap, ok := resp.Result.(map[string]any)
	assert.True(t, ok, "Result should be a map[string]any")
	assert.Equal(t, mockResult["accountBalance"], resultMap["accountBalance"])
	assert.Equal(t, mockResult["depositAddress"], resultMap["depositAddress"])
}

func TestSendRequestGraphQLError(t *testing.T) {
	operationName := "account"
	errorMessage := "Your account has been restricted."
	responseBytes := createMockGraphQLResponse(t, operationName, nil, errorMessage)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName)
	defer server.Close()

	client := dummyClient(server)

	resp, err := client.SendRequest(operationName, queries.ACCOUNT, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Equal(t, errorMessage, resp.Error)
	assert.Nil(t, resp.Result)
}

func TestSendRequestHTTPError(t *testing.T) {
	operationName := "account"
	server := mockGraphQLServer(t, []byte("Server Error"), http.StatusInternalServerError, true, operationName)
	defer server.Close()

	client := dummyClient(server)

	resp, err := client.SendRequest(operationName, queries.ACCOUNT, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Contains(t, resp.Error, "500 Internal Server Error")
	assert.Nil(t, resp.Result)
}

func TestSendRequestTypedError(t *testing.T) {
	operationName := "account"
	errorMessage := "GraphQL error occurred"
	responseBytes := createMockGraphQLResponse(t, operationName, nil, errorMessage)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true)
	defer server.Close()

	client := dummyClient(server)

	_, err := sendRequestTyped[types.Account](client, operationName, queries.ACCOUNT, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
	assert.Contains(t, err.Error(), errorMessage)
}

func TestGetAvailableCountries(t *testing.T) {
	operationName := "availableCountries"
	mockResult := []map[string]string{
		{"id": "1", "name": "Nigeria", "code": "NG"},
		{"id": "2", "name": "United States", "code": "US"},
	}
	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName)
	defer server.Close()

	client := dummyClient(server)

	countries, err := client.GetAvailableCountries()
	assert.NoError(t, err)
	assert.Len(t, countries, 2)

	assert.Equal(t, "1", countries[0].ID)
	assert.Equal(t, "Nigeria", countries[0].Name)
	assert.Equal(t, "NG", countries[0].Code)

	assert.Equal(t, "2", countries[1].ID)
	assert.Equal(t, "United States", countries[1].Name)
	assert.Equal(t, "US", countries[1].Code)
}

func TestGetMarketRate(t *testing.T) {
	operationName := "marketRate"
	countryCode := "NG"
	mockResult := map[string]float64{
		"depositRate":    1520.0,
		"withdrawalRate": 1515.0,
	}

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName, `"countryCode":"NG"`)
	defer server.Close()

	client := dummyClient(server)

	variables := map[string]string{"countryCode": countryCode}
	marketRate, err := sendRequestTyped[types.MarketRate](client, operationName, queries.MARKET_RATE, variables)

	assert.NoError(t, err)

	assert.Equal(t, 1520.0, marketRate.DepositRate)
	assert.Equal(t, 1515.0, marketRate.WithdrawalRate)
}

func TestGetPaymentMethodTypes(t *testing.T) {
	operationName := "p2pPaymentMethodTypes"
	countryID := "1"
	mockResult := []map[string]any{
		{"id": "1", "identifier": "BANK_TRANSFER", "label": "Bank Transfer"},
		{"id": "2", "identifier": "MOBILE_MONEY", "label": "Mobile Money"},
	}

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName, `"country":"1"`)
	defer server.Close()

	client := dummyClient(server)

	paymentMethodTypes, err := client.GetPaymentMethodTypes(countryID)

	assert.NoError(t, err)
	assert.Len(t, paymentMethodTypes, 2)

	assert.Equal(t, "1", paymentMethodTypes[0].ID)
	assert.Equal(t, "BANK_TRANSFER", paymentMethodTypes[0].Identifier)
	assert.Equal(t, "Bank Transfer", paymentMethodTypes[0].Label)

	assert.Equal(t, "2", paymentMethodTypes[1].ID)
	assert.Equal(t, "MOBILE_MONEY", paymentMethodTypes[1].Identifier)
	assert.Equal(t, "Mobile Money", paymentMethodTypes[1].Label)
}

func TestGetRampableAssets(t *testing.T) {
	operationName := "rampableAssets"
	mockResult := []map[string]any{
		{"name": "USDT", "symbol": "USDT", "networks": []string{"TRC20", "CELO"}},
		{"name": "USDC", "symbol": "USDC", "networks": []string{"TRC20"}},
	}

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName)
	defer server.Close()

	client := dummyClient(server)

	rampableAssets, err := client.GetRampableAssets()
	assert.NoError(t, err)
	assert.Len(t, rampableAssets, 2)

	assert.Equal(t, "USDT", rampableAssets[0].Name)
	assert.Equal(t, "USDT", rampableAssets[0].Symbol)

	assert.Equal(t, "USDC", rampableAssets[1].Name)
	assert.Equal(t, "USDC", rampableAssets[1].Symbol)
}

func TestGetRampLimits(t *testing.T) {
	operationName := "rampLimits"
	mockResult := map[string]any{
		"minimumDepositUsd": 100.0,
		"maximumDepositUsd": 10000.0,
	}

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName)
	defer server.Close()

	client := dummyClient(server)

	rampLimits, err := client.GetRampLimits()
	assert.NoError(t, err)
	assert.NotNil(t, rampLimits)

	assert.Equal(t, 100.0, rampLimits.MinimumDepositUsd)
	assert.Equal(t, 10000.0, rampLimits.MaximumDepositUsd)
}

func TestGetPaymentRequest(t *testing.T) {
	operationName := "merchantPaymentRequest"
	reference := "test_ref_1"
	mockResult := map[string]any{
		"id":          "1",
		"paymentType": "ONCHAIN",
		"hostedLink":  "https://payment-link.com",
		"amount":      100.0,
		"currency":    "USD",
		"reference":   reference,
		"status":      "PENDING",
	}

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName, `"reference":"test_ref_1"`)
	defer server.Close()

	client := dummyClient(server)

	paymentRequest, err := client.GetPaymentRequest(reference)
	assert.NoError(t, err)
	assert.NotNil(t, paymentRequest)

	assert.Equal(t, "1", paymentRequest.ID)
	assert.Equal(t, "PENDING", paymentRequest.Status)
	assert.Equal(t, reference, paymentRequest.Reference)
}

func TestGetAccount(t *testing.T) {
	operationName := "account"
	mockResult := map[string]any{
		"id":             "1",
		"accountBalance": 100.0,
		"depositAddress": "deposit-address",
	}

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName)
	defer server.Close()

	client := dummyClient(server)

	account, err := client.GetAccount()
	assert.NoError(t, err)
	assert.NotNil(t, account)

	assert.Equal(t, "1", account.ID)
	assert.Equal(t, 100.0, account.AccountBalance)
	assert.Equal(t, "deposit-address", account.DepositAddress)
}

func TestConfirmTransaction(t *testing.T) {
	operationName := "confirmTransaction"
	mockResult := true

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName, `"paymentRequest":"1"`, `"transactionHash":"tx_hash"`)
	defer server.Close()

	client := dummyClient(server)

	confirmTransactionInput := types.ConfirmTransactionInput{
		PaymentRequest:  "1",
		TransactionHash: "tx_hash",
	}

	confirmed, err := client.ConfirmTransaction(confirmTransactionInput)
	assert.NoError(t, err)
	assert.True(t, confirmed)
}

func TestInitiateHostedPayment(t *testing.T) {
	operationName := "initiateHostedPayment"
	mockResult := map[string]any{
		"id":         "1",
		"hostedLink": "https://payment-link.com",
	}

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName, `"reference":"ref123"`)
	defer server.Close()

	client := dummyClient(server)

	initiateHostedPaymentInput := types.InitiateHostedPaymentInput{
		PaymentType: "ONCHAIN",
		Amount:      100.0,
		Currency:    "USD",
		Reference:   "ref123",
		RedirectUrl: "https://redirect.com",
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
	}
	initiatedPayment, err := client.InitiateHostedPayment(initiateHostedPaymentInput)
	assert.NoError(t, err)
	assert.NotNil(t, initiatedPayment)
	assert.Equal(t, "1", initiatedPayment.Id)
	assert.Equal(t, "https://payment-link.com", initiatedPayment.HostedLink)
}

func TestCancelHostedPayment(t *testing.T) {
	operationName := "cancelHostedPayment"
	mockResult := true

	responseBytes := createMockGraphQLResponse(t, operationName, mockResult)
	server := mockGraphQLServer(t, responseBytes, http.StatusOK, true, operationName, `"paymentRequest":"1"`)
	defer server.Close()

	client := dummyClient(server)

	input := types.CancelHostedPaymentInput{PaymentRequest: "1"}
	cancelled, err := client.CancelHostedPayment(input)

	assert.NoError(t, err)
	assert.True(t, cancelled)
}
