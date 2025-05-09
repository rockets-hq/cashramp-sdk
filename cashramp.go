package cashrampsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/rockets-hq/cashramp-sdk/mutations"
	"github.com/rockets-hq/cashramp-sdk/queries"
	"github.com/rockets-hq/cashramp-sdk/types"
)

const host = "api.useaccrue.com"

type Client struct {
	ApiUrl     string
	secretKey  string
	httpClient *http.Client
}

type CashrampResponse struct {
	Success bool `json:"success"`
	Result  any
	Error   string
}

type reqBody struct {
	Query     string `json:"query"`
	Variables any    `json:"variables"`
}

type graphqlErrorResponse struct {
	Message string `json:"message"`
}

type rawGraphQLResponse struct {
	Data   map[string]any         `json:"data"`
	Errors []graphqlErrorResponse `json:"errors"`
}

func InitialiseClient(environment, secretKey string) (*Client, error) {
	apiUrl, err := validateEnv(environment)
	if err != nil {
		return nil, err
	}

	secret, err := validateSecretKey(secretKey)
	if err != nil {
		return nil, err
	}

	return &Client{
		ApiUrl:     apiUrl,
		secretKey:  secret,
		httpClient: http.DefaultClient,
	}, nil
}

func (c *Client) SendRequest(name, query string, variables any) (*CashrampResponse, error) {
	response := &CashrampResponse{}
	requestBody := &reqBody{
		Query:     query,
		Variables: variables,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.ApiUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.secretKey))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		graphqlResponse := &rawGraphQLResponse{}
		jsonErr := json.NewDecoder(resp.Body).Decode(graphqlResponse)
		if jsonErr != nil {
			response.Success = false
			response.Error = jsonErr.Error()
			return response, jsonErr
		}

		if graphqlResponse.Errors != nil {
			response.Success = false
			response.Error = graphqlResponse.Errors[0].Message
			return response, nil
		} else {
			response.Success = true
			response.Result = graphqlResponse.Data[name]
		}
	default:
		response.Success = false
		response.Error = resp.Status
		return response, nil
	}

	return response, nil
}

func (c *Client) GetAvailableCountries() ([]types.Country, error) {
	return SendRequestTyped[[]types.Country](c, "availableCountries", queries.AVAILABLE_COUNTRIES, nil)
}

func (c *Client) GetMarketRate(countryCode string) (*types.MarketRate, error) {
	variables := map[string]string{
		"countryCode": countryCode,
	}

	marketRate, err := SendRequestTyped[types.MarketRate](c, "marketRate", queries.MARKET_RATE, variables)
	if err != nil {
		return nil, err
	}
	return &marketRate, nil
}

func (c *Client) GetPaymentMethodTypes(countryId string) ([]types.PaymentMethodTypes, error) {
	variables := map[string]string{
		"country": countryId,
	}

	paymentMethodTypes, err := SendRequestTyped[[]types.PaymentMethodTypes](c, "p2pPaymentMethodTypes", queries.PAYMENT_METHOD_TYPES, variables)
	if err != nil {
		return nil, err
	}
	return paymentMethodTypes, nil
}

func (c *Client) GetRampableAssets() ([]types.RampableAssets, error) {
	rampableAssets, err := SendRequestTyped[[]types.RampableAssets](c, "rampableAssets", queries.RAMPABLE_ASSETS, nil)
	if err != nil {
		return nil, err
	}

	return rampableAssets, nil
}

func (c *Client) GetRampLimits() (*types.RampLimits, error) {
	rampLimits, err := SendRequestTyped[types.RampLimits](c, "rampLimits", queries.RAMP_LIMITS, nil)
	if err != nil {
		return nil, err
	}

	return &rampLimits, err
}

func (c *Client) GetPaymentRequest(reference string) (*types.PaymentRequest, error) {
	variables := map[string]string{
		"reference": reference,
	}
	paymentRequest, err := SendRequestTyped[types.PaymentRequest](c, "merchantPaymentRequest", queries.PAYMENT_REQUEST, variables)
	if err != nil {
		return nil, err
	}

	return &paymentRequest, nil
}

func (c *Client) GetAccount() (*types.Account, error) {
	account, err := SendRequestTyped[types.Account](c, "account", queries.ACCOUNT, nil)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// Mutations

func (c *Client) ConfirmTransaction(paymentRequest types.ConfirmTransactionInput) (bool, error) {
	confirmedPayment, err := SendRequestTyped[bool](c, "confirmTransaction", mutations.CONFIRM_TRANSACTION, paymentRequest)
	if err != nil {
		return false, err
	}
	return confirmedPayment, nil
}

func (c *Client) InitiateHostedPayment(payment types.InitiateHostedPaymentInput) (*types.HostedPaymentResponse, error) {
	initiatedPayment, err := SendRequestTyped[types.HostedPaymentResponse](c, "initiateHostedPayment", mutations.INITIATE_HOSTED_PAYMENT, payment)
	if err != nil {
		return nil, err
	}
	return &initiatedPayment, nil
}

func (c *Client) CancelHostedPayment(payment types.CancelHostedPaymentInput) (bool, error) {
	initiatedPayment, err := SendRequestTyped[bool](c, "cancelHostedPayment", mutations.CANCEL_HOSTED_PAYMENT, payment)
	if err != nil {
		return false, err
	}
	return initiatedPayment, nil
}

func (c *Client) CreateCustomer(customer types.CreateCustomerInput) (*types.Customer, error) {
	createdCustomer, err := SendRequestTyped[types.Customer](c, "createCustomer", mutations.CREATE_CUSTOMER, customer)
	if err != nil {
		return nil, err
	}
	return &createdCustomer, nil
}

func (c *Client) AddPaymentMethod(payment types.AddPaymentMethodInput) (*types.AddPaymentMethodResponse, error) {
	initiatedPayment, err := SendRequestTyped[types.AddPaymentMethodResponse](c, "addPaymentMethod", mutations.ADD_PAYMENT_METHOD, payment)
	if err != nil {
		return nil, err
	}
	return &initiatedPayment, nil
}

func (c *Client) WithdrawOnchain(payment types.WithdrawOnchainInput) (*types.WithdrawOnchainResponse, error) {
	initiatedPayment, err := SendRequestTyped[types.WithdrawOnchainResponse](c, "withdrawOnchain", mutations.WITHDRAW_ONCHAIN, payment)
	if err != nil {
		return nil, err
	}
	return &initiatedPayment, nil
}

// TODO: return error message from the server when there is one
func SendRequestTyped[T any](client *Client, name, query string, variables any) (T, error) {
	var out T
	resp, err := client.SendRequest(name, query, variables)
	if err != nil {
		return out, err
	}

	if !resp.Success {
		return out, fmt.Errorf("request failed: %s", resp.Error)
	}

	// Convert the generic result into typed output
	raw, err := json.Marshal(resp.Result)
	if err != nil {
		return out, err
	}
	err = json.Unmarshal(raw, &out)
	return out, err
}

func validateEnv(env string) (apiUrl string, err error) {
	var environment string
	if env == "" {
		environment = os.Getenv("CASHRAMP_ENV")
	} else {
		environment = env
	}

	switch environment {
	case "test":
		apiUrl = fmt.Sprintf("https://staging.%v/cashramp/api/graphql", host)
	case "live":
		apiUrl = fmt.Sprintf("https://%v/cashramp/api/graphql", host)
	default:
		err = fmt.Errorf(`%v is not a valid env. Can either be "test" or "live"`, environment)
	}
	return
}

func validateSecretKey(secretKey string) (string, error) {
	var secret string

	switch secretKey {
	case "":
		secretFromEnv := os.Getenv("CASHRAMP_SECRET_KEY")
		if secretFromEnv == "" {
			return "", errors.New("please provide your API secret key")
		} else {
			secret = secretFromEnv
		}
	default:
		secret = secretKey
	}

	return secret, nil
}
