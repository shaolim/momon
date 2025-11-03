package model

import "encoding/json"

type Receipt struct {
	Shop            string  `json:"shop"`
	TransactionDate string  `json:"transactionDate"`
	Items           []Item  `json:"items"`
	Tax             float64 `json:"tax"`
	Total           float64 `json:"total"`
	IsValid         bool    `json:"isValid"`
	Message         string  `json:"message"`
}

func (r *Receipt) String() string {
	jsonBytes, _ := json.Marshal(r)
	return string(jsonBytes)
}

// TODO: validate the calculation

type Item struct {
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	TotalPrice float64 `json:"totalPrice"`
}
