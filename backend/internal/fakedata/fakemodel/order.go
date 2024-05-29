package fakemodel

import "time"

type OrderFake struct {
	OrderUID          string
	TrackNumber       string
	Entry             string
	Delivery          DeliveryFake
	Payment           PaymentFake
	Items             []ItemFake
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	Shardkey          string
	SMID              int
	DateCreated       time.Time
	OofShard          string
}

type DeliveryFake struct {
	Name    string
	Phone   string `fake:"{phone}"`
	Zip     string
	City    string
	Address string
	Region  string
	Email   string `fake:"{email}"`
}

type ItemFake struct {
	ChrtID      int `fake:"{number:10}"`
	TrackNumber string
	Price       int `fake:"{number:10}"`
	Rid         string
	Name        string
	Sale        int `fake:"{number:10}"`
	Size        string
	TotalPrice  int `fake:"{number:10}"`
	NMID        int `fake:"{number:10}"`
	Brand       string
	Status      int `fake:"{number:10}"`
}

type PaymentFake struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int `fake:"{number:10}"`
	PaymentDt    int `fake:"{number:10}"`
	Bank         string
	DeliveryCost int `fake:"{number:10}"`
	GoodsTotal   int `fake:"{number:10}"`
	CustomFee    int `fake:"{number:10}"`
}
