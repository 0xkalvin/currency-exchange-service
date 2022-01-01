package entities

type Order struct {
	Id               string `json:"id,omitempty"`
	Amount           string `json:"amount"`
	CustomerId       string `json:"customer_id"`
	SourceCurrencyId string `json:"source_currency_id"`
	TargetCurrencyId string `json:"target_currency_id"`
}
