package main

// Transaction struct
type Transaction struct {
	State   string  `json:"state"`
	Amount  float64 `json:"amount"`
	ID      int     `json:"transaction_id"`
	SrcType string
}
