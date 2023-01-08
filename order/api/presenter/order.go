package presenter

type Order struct {
	WalletID  int64   `json:"walletId"`
	OrderType string  `json:"type"`
	Asset     string  `json:"asset"`
	Quantity  float64 `json:"quantity"`
}
