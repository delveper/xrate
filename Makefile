ENV := .env
include $(ENV)

rate-local:
	curl https://api.coingecko.com/api/v3/exchange_rates | jq '.rates.uah.value'