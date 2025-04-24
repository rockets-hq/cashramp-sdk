<p align="center">
  <img alt="Last Commit" src="https://badgen.net/github/last-commit/rockets-hq/cashramp-sdk" />
  <a href="https://github.com/rockets-hq/cashramp-sdk/"><img src="https://img.shields.io/github/stars/rockets-hq/cashramp-sdk.svg"/></a>
  <a href="https://github.com/rockets-hq/cashramp-sdk/"><img src="https://img.shields.io/npm/l/cashramp.svg"/></a>
</p>

# Cashramp SDK Go

This is the official Go SDK for [Cashramp's API](https://cashramp.co/commerce).

### ➕ Installation

```bash
go get github.com/rockets-hq/cashramp-sdk-go
```

### 👨🏾‍💻 Quick Start

```go
cashrampApi, err := cashrampsdk.InitialiseClient(
	"test", //can be either test or live
	"CSHRMP-SECK_apE0rjq1tiWl6VLB",
)
if err != nil {
	panic(err)
}
	
// Example: Fetch available countries
countries, err := cashrampApi.GetAvailableCountries()
if err != nil {
	log.Println(err)
}

log.Println(countries)
```

## API Reference

### Queries

- `getAvailableCountries()`: Fetch the countries that Cashramp is available in
- `getMarketRate({ countryCode })`: Fetch the Cashramp market rate for a country
- `getPaymentMethodTypes({ country })`: Fetch the payment method types available in a country
- `getRampableAssets()`: Fetch the assets you can on/offramp with the Onchain Ramp
- `getRampLimits()`: Fetch the Onchain Ramp limits
- `getPaymentRequest({ reference })`: Fetch the details of a payment request
- `getAccount()`: Fetch the account information for the authenticated user.

### Mutations

- `confirmTransaction({ paymentRequest, transactionHash })`: Confirm a crypto transfer sent into Cashramp's Secure Escrow address
- `initiateHostedPayment({ amount, paymentType, countryCode, currency, email, reference, redirectUrl, firstName, lastName })`: Initiate a payment request
- `cancelHostedPayment({ paymentRequest })`: Cancel an ongoing payment request
- `createCustomer({ firstName, lastName, email, country })`: Create a new customer profile
- `addPaymentMethod({ customer, paymentMethodType, fields })`: Add a payment method for an existing customer
- `withdrawOnchain({ address, amountUsd })`:  Withdraw from your balance to an onchain wallet address


## Custom Queries

For advanced use cases where the provided methods don't cover your specific needs, you can use the `sendRequest` method to send custom GraphQL queries:

```go
query := `query {
			availableCountries {
				id
				name
				code
				currency {
					isoCode
					name
				}
			}
		}`

response, err := cashrampApi.SendRequest("availableCountries", query, nil)
if err != nil {
	log.Println(err)
}

if !response.Success {
	fmt.Printf("request failed: %s", response.Error)
}

fmt.Printf("response result: %v", response.Result)
```

## Error Handling

All methods in the SDK return an error value `err` which will contain details about the error. For more complex queries where `SendRequest` is used, the response object contains a `success` boolean. When `success` is `false`, an `Error` field will be available with details about the error.

## Go Support

This SDK includes Go struct types out of the box.

## Documentation

For detailed API documentation, visit [Cashramp's API docs](https://docs.cashramp.co).