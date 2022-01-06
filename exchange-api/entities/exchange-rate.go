package entities

type ExchangeRate struct {
	Id             string `json:"id,omitempty"`
	CurrencyId     string `json:"currency_id"`
	BaseCurrencyId string `json:"base_currency_id"`
	Rate           string `json:"rate"`
	Timestamp      string `json:"timestamp"`
}
