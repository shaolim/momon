package model_test

import (
	"encoding/json"
	"github/shaolim/momon/internal/receipt/model"
	"testing"
)

func TestReceipt_String(t *testing.T) {
	tests := []struct {
		name    string
		receipt model.Receipt
		want    string
	}{
		{
			name: "valid receipt with single item",
			receipt: model.Receipt{
				Shop:            "Test Store",
				TransactionDate: "2024-01-15 14:30",
				Items: []model.Item{
					{
						Name:       "Apple",
						Quantity:   2,
						Price:      100,
						Tax:        10,
						TotalPrice: 210,
					},
				},
				Tax:     10,
				Total:   210,
				IsValid: true,
				Message: "",
			},
			want: `{"shop":"Test Store","transactionDate":"2024-01-15 14:30","items":[{"name":"Apple","quantity":2,"price":100,"tax":10,"totalPrice":210}],"tax":10,"total":210,"isValid":true,"message":""}`,
		},
		{
			name: "valid receipt with multiple items",
			receipt: model.Receipt{
				Shop:            "Grocery Market",
				TransactionDate: "2024-02-20 09:15",
				Items: []model.Item{
					{
						Name:       "Bread",
						Quantity:   1,
						Price:      50,
						Tax:        5,
						TotalPrice: 55,
					},
					{
						Name:       "Milk",
						Quantity:   2,
						Price:      80,
						Tax:        16,
						TotalPrice: 176,
					},
				},
				Tax:     21,
				Total:   231,
				IsValid: true,
				Message: "",
			},
			want: `{"shop":"Grocery Market","transactionDate":"2024-02-20 09:15","items":[{"name":"Bread","quantity":1,"price":50,"tax":5,"totalPrice":55},{"name":"Milk","quantity":2,"price":80,"tax":16,"totalPrice":176}],"tax":21,"total":231,"isValid":true,"message":""}`,
		},
		{
			name: "empty receipt",
			receipt: model.Receipt{
				Shop:            "",
				TransactionDate: "",
				Items:           []model.Item{},
				Tax:             0,
				Total:           0,
				IsValid:         false,
				Message:         "",
			},
			want: `{"shop":"","transactionDate":"","items":[],"tax":0,"total":0,"isValid":false,"message":""}`,
		},
		{
			name: "receipt with nil items",
			receipt: model.Receipt{
				Shop:            "Test Shop",
				TransactionDate: "2024-03-10 12:00",
				Items:           nil,
				Tax:             0,
				Total:           0,
				IsValid:         true,
				Message:         "",
			},
			want: `{"shop":"Test Shop","transactionDate":"2024-03-10 12:00","items":null,"tax":0,"total":0,"isValid":true,"message":""}`,
		},
		{
			name: "receipt with zero values",
			receipt: model.Receipt{
				Shop:            "Zero Store",
				TransactionDate: "2024-04-05 00:00",
				Items: []model.Item{
					{
						Name:       "Free Item",
						Quantity:   0,
						Price:      0,
						Tax:        0,
						TotalPrice: 0,
					},
				},
				Tax:     0,
				Total:   0,
				IsValid: true,
				Message: "",
			},
			want: `{"shop":"Zero Store","transactionDate":"2024-04-05 00:00","items":[{"name":"Free Item","quantity":0,"price":0,"tax":0,"totalPrice":0}],"tax":0,"total":0,"isValid":true,"message":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.receipt.String()
			if got != tt.want {
				t.Errorf("Receipt.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReceipt_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		receipt model.Receipt
		wantErr bool
	}{
		{
			name: "marshal and unmarshal receipt with items",
			receipt: model.Receipt{
				Shop:            "Test Store",
				TransactionDate: "2024-01-15 14:30",
				Items: []model.Item{
					{
						Name:       "Laptop",
						Quantity:   1,
						Price:      1000,
						Tax:        100,
						TotalPrice: 1100,
					},
				},
				Tax:     100,
				Total:   1100,
				IsValid: true,
				Message: "",
			},
			wantErr: false,
		},
		{
			name: "marshal empty receipt",
			receipt: model.Receipt{
				Shop:            "",
				TransactionDate: "",
				Items:           []model.Item{},
				Tax:             0,
				Total:           0,
				IsValid:         false,
				Message:         "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			jsonBytes, err := json.Marshal(tt.receipt)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Unmarshal
			var unmarshaled model.Receipt
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare
			if unmarshaled.Shop != tt.receipt.Shop {
				t.Errorf("Shop = %v, want %v", unmarshaled.Shop, tt.receipt.Shop)
			}
			if unmarshaled.TransactionDate != tt.receipt.TransactionDate {
				t.Errorf("TransactionDate = %v, want %v", unmarshaled.TransactionDate, tt.receipt.TransactionDate)
			}
			if unmarshaled.Tax != tt.receipt.Tax {
				t.Errorf("Tax = %v, want %v", unmarshaled.Tax, tt.receipt.Tax)
			}
			if unmarshaled.Total != tt.receipt.Total {
				t.Errorf("Total = %v, want %v", unmarshaled.Total, tt.receipt.Total)
			}
			if unmarshaled.IsValid != tt.receipt.IsValid {
				t.Errorf("IsValid = %v, want %v", unmarshaled.IsValid, tt.receipt.IsValid)
			}
			if len(unmarshaled.Items) != len(tt.receipt.Items) {
				t.Errorf("Items length = %v, want %v", len(unmarshaled.Items), len(tt.receipt.Items))
			}
		})
	}
}

func TestItem_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		item    model.Item
		wantErr bool
	}{
		{
			name: "marshal and unmarshal item",
			item: model.Item{
				Name:       "Coffee",
				Quantity:   3,
				Price:      250,
				Tax:        75,
				TotalPrice: 825,
			},
			wantErr: false,
		},
		{
			name: "marshal item with zero values",
			item: model.Item{
				Name:       "Free Sample",
				Quantity:   0,
				Price:      0,
				Tax:        0,
				TotalPrice: 0,
			},
			wantErr: false,
		},
		{
			name: "marshal item with decimal values",
			item: model.Item{
				Name:       "Banana",
				Quantity:   1.5,
				Price:      99.99,
				Tax:        15.99,
				TotalPrice: 165.98,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			jsonBytes, err := json.Marshal(tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Unmarshal
			var unmarshaled model.Item
			err = json.Unmarshal(jsonBytes, &unmarshaled)
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare
			if unmarshaled.Name != tt.item.Name {
				t.Errorf("Name = %v, want %v", unmarshaled.Name, tt.item.Name)
			}
			if unmarshaled.Quantity != tt.item.Quantity {
				t.Errorf("Quantity = %v, want %v", unmarshaled.Quantity, tt.item.Quantity)
			}
			if unmarshaled.Price != tt.item.Price {
				t.Errorf("Price = %v, want %v", unmarshaled.Price, tt.item.Price)
			}
			if unmarshaled.Tax != tt.item.Tax {
				t.Errorf("Tax = %v, want %v", unmarshaled.Tax, tt.item.Tax)
			}
			if unmarshaled.TotalPrice != tt.item.TotalPrice {
				t.Errorf("TotalPrice = %v, want %v", unmarshaled.TotalPrice, tt.item.TotalPrice)
			}
		})
	}
}

func TestReceipt_StructCreation(t *testing.T) {
	t.Run("create receipt with all fields", func(t *testing.T) {
		receipt := model.Receipt{
			Shop:            "My Shop",
			TransactionDate: "2024-05-01 10:00",
			Items: []model.Item{
				{Name: "Item1", Quantity: 1, Price: 100, Tax: 10, TotalPrice: 110},
			},
			Tax:     10,
			Total:   110,
			IsValid: true,
			Message: "",
		}

		if receipt.Shop != "My Shop" {
			t.Errorf("Shop = %v, want My Shop", receipt.Shop)
		}
		if receipt.TransactionDate != "2024-05-01 10:00" {
			t.Errorf("TransactionDate = %v, want 2024-05-01 10:00", receipt.TransactionDate)
		}
		if receipt.Tax != 10 {
			t.Errorf("Tax = %v, want 10", receipt.Tax)
		}
		if receipt.Total != 110 {
			t.Errorf("Total = %v, want 110", receipt.Total)
		}
		if !receipt.IsValid {
			t.Error("IsValid = false, want true")
		}
		if len(receipt.Items) != 1 {
			t.Errorf("Items length = %v, want 1", len(receipt.Items))
		}
	})

	t.Run("create receipt with default values", func(t *testing.T) {
		receipt := model.Receipt{}

		if receipt.Shop != "" {
			t.Errorf("Shop = %v, want empty string", receipt.Shop)
		}
		if receipt.Tax != 0 {
			t.Errorf("Tax = %v, want 0", receipt.Tax)
		}
		if receipt.Total != 0 {
			t.Errorf("Total = %v, want 0", receipt.Total)
		}
		if receipt.IsValid {
			t.Error("IsValid = true, want false")
		}
		if receipt.Items != nil {
			t.Errorf("Items = %v, want nil", receipt.Items)
		}
	})
}

func TestItem_StructCreation(t *testing.T) {
	t.Run("create item with all fields", func(t *testing.T) {
		item := model.Item{
			Name:       "Test Item",
			Quantity:   5,
			Price:      200,
			Tax:        50,
			TotalPrice: 1050,
		}

		if item.Name != "Test Item" {
			t.Errorf("Name = %v, want Test Item", item.Name)
		}
		if item.Quantity != 5 {
			t.Errorf("Quantity = %v, want 5", item.Quantity)
		}
		if item.Price != 200 {
			t.Errorf("Price = %v, want 200", item.Price)
		}
		if item.Tax != 50 {
			t.Errorf("Tax = %v, want 50", item.Tax)
		}
		if item.TotalPrice != 1050 {
			t.Errorf("TotalPrice = %v, want 1050", item.TotalPrice)
		}
	})

	t.Run("create item with default values", func(t *testing.T) {
		item := model.Item{}

		if item.Name != "" {
			t.Errorf("Name = %v, want empty string", item.Name)
		}
		if item.Quantity != 0 {
			t.Errorf("Quantity = %v, want 0", item.Quantity)
		}
		if item.Price != 0 {
			t.Errorf("Price = %v, want 0", item.Price)
		}
		if item.Tax != 0 {
			t.Errorf("Tax = %v, want 0", item.Tax)
		}
		if item.TotalPrice != 0 {
			t.Errorf("TotalPrice = %v, want 0", item.TotalPrice)
		}
	})
}

func TestReceipt_String_EdgeCases(t *testing.T) {
	t.Run("receipt with special characters in shop name", func(t *testing.T) {
		receipt := model.Receipt{
			Shop:            "Test & Co. \"Shop\"",
			TransactionDate: "2024-01-01 00:00",
			Items:           []model.Item{},
			Tax:             0,
			Total:           0,
			IsValid:         true,
			Message:         "",
		}

		result := receipt.String()
		if result == "" {
			t.Error("String() returned empty string")
		}

		// Verify it's valid JSON
		var unmarshaled model.Receipt
		err := json.Unmarshal([]byte(result), &unmarshaled)
		if err != nil {
			t.Errorf("String() returned invalid JSON: %v", err)
		}

		if unmarshaled.Shop != receipt.Shop {
			t.Errorf("Shop after unmarshal = %v, want %v", unmarshaled.Shop, receipt.Shop)
		}
	})

	t.Run("receipt with very large numbers", func(t *testing.T) {
		receipt := model.Receipt{
			Shop:            "Big Numbers Store",
			TransactionDate: "2024-06-01 12:00",
			Items: []model.Item{
				{
					Name:       "Expensive Item",
					Quantity:   1000000,
					Price:      999999.99,
					Tax:        99999.99,
					TotalPrice: 1099999989.99,
				},
			},
			Tax:     99999.99,
			Total:   1099999989.99,
			IsValid: true,
			Message: "",
		}

		result := receipt.String()
		if result == "" {
			t.Error("String() returned empty string")
		}

		// Verify it's valid JSON
		var unmarshaled model.Receipt
		err := json.Unmarshal([]byte(result), &unmarshaled)
		if err != nil {
			t.Errorf("String() returned invalid JSON: %v", err)
		}
	})

	t.Run("receipt with negative values", func(t *testing.T) {
		receipt := model.Receipt{
			Shop:            "Refund Store",
			TransactionDate: "2024-07-01 15:30",
			Items: []model.Item{
				{
					Name:       "Returned Item",
					Quantity:   -1,
					Price:      100,
					Tax:        -10,
					TotalPrice: -110,
				},
			},
			Tax:     -10,
			Total:   -110,
			IsValid: true,
			Message: "",
		}

		result := receipt.String()
		if result == "" {
			t.Error("String() returned empty string")
		}

		// Verify it's valid JSON
		var unmarshaled model.Receipt
		err := json.Unmarshal([]byte(result), &unmarshaled)
		if err != nil {
			t.Errorf("String() returned invalid JSON: %v", err)
		}

		if unmarshaled.Total != receipt.Total {
			t.Errorf("Total after unmarshal = %v, want %v", unmarshaled.Total, receipt.Total)
		}
	})
}
