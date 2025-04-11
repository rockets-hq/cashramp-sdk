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
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Networks        string `json:"networks"`
	ContractAddress string `json:"contractAddress"`
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

type HostedPayment struct {
	PaymentType string  `json:"paymentType"`
	Amount      float64 `json:"amount"`
	CountryCode string  `json:"countryCode"`
	Reference   string  `json:"reference"`
	RedirectUrl string  `json:"redirectUrl"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
}

type HostedPaymentResponse struct {
	Id         string        `json:"id"`
	HostedLink string        `json:"hostedLink"`
	Status     PaymentStatus `json:"status"`
}

type NewCustomer struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CountryID   string `json:"country"`
}

type Customer struct {
	Id        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Country   Country `json:"country"`
}
