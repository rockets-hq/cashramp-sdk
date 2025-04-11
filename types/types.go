package types

type PaymentStatus string

const (
	PaymentStatusTypeCreated    PaymentStatus = "created"
	PaymnetStatusTypedPickedUp  PaymentStatus = "picked_up"
	PaymentStatusTypedCompleted PaymentStatus = "completed"
	PaymentStatusTypedCancelled PaymentStatus = "cancelled"
)

type Country struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type MarketRate struct {
	DepositRate    float64 `json:"depositRate"`
	WithdrawalRate float64 `json:"withdrawalRate"`
}

type PaymentMethodTypes struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	Label      string `json:"label"`
	Fields     struct {
		Label      string `json:"label"`
		Identifier string `json:"identifier"`
		Required   bool   `json:"required"`
	}
}

type RampableAssets struct {
	Name            string   `json:"name"`
	Symbol          string   `json:"symbol"`
	Networks        []string `json:"networks"`
	ContractAddress string   `json:"contractAddress"`
}

type RampLimits struct {
	MinimumDepositUsd    float64 `json:"minimumDepositUsd"`
	MaximumDepositUsd    float64 `json:"maximumDepositUsd"`
	MinimumWithdrawalUsd float64 `json:"minimumWithdrawalUsd"`
	MaximumWithdrawalUsd float64 `json:"maximumWithdrawalUsd"`
	DailyLimitUsd        float64 `json:"dailyLimitUsd"`
}

type PaymentRequest struct {
	ID          string  `json:"id"`
	PaymentType string  `json:"paymentType"`
	HostedLink  string  `json:"hostedLink"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Reference   string  `json:"reference"`
	Status      string  `json:"status"`
}

type Account struct {
	ID             string  `json:"id"`
	AccountBalance float64 `json:"accountBalance"`
	DepositAddress string  `json:"depositAddress"`
}

type ConfirmTransactionInput struct {
	PaymentRequest  string `json:"paymentRequest"`
	TransactionHash string `json:"transactionHash"`
}

type InitiateHostedPaymentInput struct {
	PaymentType string  `json:"paymentType"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	CountryCode string  `json:"countryCode"`
	Reference   string  `json:"reference"`
	RedirectUrl string  `json:"redirectUrl"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
}

type CancelHostedPaymentInput struct {
	PaymentRequest string `json:"paymentRequest"`
}

type HostedPaymentResponse struct {
	Id         string        `json:"id"`
	HostedLink string        `json:"hostedLink"`
	Status     PaymentStatus `json:"status"`
}

type CreateCustomerInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CountryID string `json:"country"`
}

type Customer struct {
	Id        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Country   Country `json:"country"`
}

type AddPaymentMethodInput struct {
	CustomerID          string `json:"customer"`
	PaymentMethodTypeID string `json:"p2pPaymentMethodType"`
	Fields              []struct {
		Identifier string `json:"identifier"`
		Value      string `json:"value"`
	} `json:"fields"`
}

type AddPaymentMethodResponse struct {
	ID     string `json:"id"`
	Value  string `json:"value"`
	Fields []struct {
		Identifier string `json:"identifier"`
		Value      string `json:"value"`
	} `json:"fields"`
}
type WithdrawOnchainInput struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type WithdrawOnchainResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
