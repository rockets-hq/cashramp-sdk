package cashrampsdkgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/rockets-hq/cashramp-sdk-go/queries"
)

const host = "api.useaccrue.com"

type Client struct {
	env       string
	apiUrl    string
	secretKey string

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
	env, apiUrl, err := validateEnv(environment)
	if err != nil {
		return nil, err
	}

	secret, err := validateSecretKey(secretKey)
	if err != nil {
		return nil, err
	}

	return &Client{
		env:        env,
		apiUrl:     apiUrl,
		secretKey:  secret,
		httpClient: http.DefaultClient,
	}, nil
}

func (c *Client) GetAvailableCountries() (*CashrampResponse, error) {
	return c.SendRequest("availableCountries", queries.AVAILABLE_COUNTRIES, nil)
}

func (c *Client) GetMarketRate(countryCode string) (*CashrampResponse, error) {
	variables := map[string]string{
		"countryCode": countryCode,
	}

	return c.SendRequest("marketRate", queries.MARKET_RATE, variables)
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

	req, err := http.NewRequest(http.MethodPost, c.apiUrl, bytes.NewBuffer(body))
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

func validateEnv(env string) (environment string, apiUrl string, err error) {
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
			return "", errors.New("Please provide your API secret key")
		} else {
			secret = secretFromEnv
		}
	default:
		secret = secretKey
	}

	return secret, nil
}
