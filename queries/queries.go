package queries

const (
	AVAILABLE_COUNTRIES = `
		query {
			availableCountries {
				id
				name
				code
			}
		}
	`

	MARKET_RATE = `
		query ($countryCode: String!) {
				marketRate(countryCode: $countryCode) {
					depositRate
					withdrawalRate
			}
		}
	`

	PAYMENT_METHOD_TYPES = `
		query ($country: ID!) {
			p2pPaymentMethodTypes(country: $country) {
				id
				identifier
				label
				fields {
					label
					identifier
					required
				}
			}
		}
	`

	RAMPABLE_ASSETS = `
		query {
			rampableAssets {
				name
				symbol
				networks
				contractAddress
			}
		}
	`

	RAMP_LIMITS = `
		query {
			rampLimits {
				minimumDepositUsd
				maximumDepositUsd
				minimumWithdrawalUsd
				maximumWithdrawalUsd
				dailyLimitUsd
			}
		}
	`
	
	PAYMENT_REQUEST = `
		query ($reference: String!) {
			merchantPaymentRequest(reference: $reference) {
				id
				paymentType
				hostedLink
				amount
				currency
				reference
				status
			}
		}
	`
	
	ACCOUNT = `
		query {
			account {
				id
				accountBalance
				depositAddress
			}
		}
	`
)

