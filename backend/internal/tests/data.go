package tests

import (
	"time"
	"wb/backend/internal/domain/entity"
)

var (
	order1 entity.Order = entity.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: entity.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com"},
		Payment: entity.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0},
		Items: []entity.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NMID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SMID:              99,
		DateCreated:       func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-11-26T06:17:33Z"); return t }(),
		OofShard:          "1",
	}

	order2 entity.Order = entity.Order{
		OrderUID:    "order2",
		TrackNumber: "order2",
		Entry:       "order2WBIL",
		Delivery: entity.Delivery{
			Name:    "order2",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "order2",
			Region:  "order2",
			Email:   "test@gmail.com"},
		Payment: entity.Payment{
			Transaction:  "order2",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0},
		Items: []entity.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "order2",
				Price:       453,
				Rid:         "order2",
				Name:        "order2",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NMID:        2389212,
				Brand:       "order2 Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SMID:              99,
		DateCreated:       func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-11-26T06:17:33Z"); return t }(),
		OofShard:          "1",
	}
)
