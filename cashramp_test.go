package cashrampsdkgo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockGraphQLServer(t *testing.T, mockResponse []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

		w.WriteHeader(http.StatusOK)
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

func TestGetAvailableCountries(t *testing.T) {
	mockGraphQLResponse := map[string]any{
		"data": map[string]any{
			"availableCountries": []map[string]string{
				{"id": "1", "name": "Nigeria", "code": "NG"},
				{"id": "2", "name": "United States", "code": "US"},
			},
		},
	}

	responseBytes, _ := json.Marshal(mockGraphQLResponse)

	// Mock GraphQL server
	server := mockGraphQLServer(t, responseBytes)
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
	mockGraphQLResponse := map[string]any{
		"data": map[string]any{
			"marketRate": map[string]float64{
				"depositRate":    1520.0,
				"withdrawalRate": 1515.0,
			},
		},
	}

	responseBytes, _ := json.Marshal(mockGraphQLResponse)
	server := mockGraphQLServer(t, responseBytes)
	defer server.Close()

	client := dummyClient(server)

	marketRate, err := client.GetMarketRate("NG")
	assert.NoError(t, err)
	assert.NotNil(t, marketRate)

	assert.Equal(t, 1520.0, marketRate.DepositRate)
	assert.Equal(t, 1515.0, marketRate.WithdrawalRate)
}
