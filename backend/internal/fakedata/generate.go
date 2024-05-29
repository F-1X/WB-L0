package fakedata

import (
	"log"
	"wb/backend/internal/domain/entity"
	"wb/backend/internal/fakedata/fakemodel"

	"github.com/brianvoe/gofakeit/v7"
)

// Generate - генерация рандомных данных (см fakemodel теги)
func Generate(count int) []entity.Order {
	var fakes []entity.Order
	for i := 0; i < count; i++ {
		var f fakemodel.OrderFake
		err := gofakeit.Struct(&f)
		if err != nil {
			log.Fatal("something broked")
		}

		order := entity.Order{
			OrderUID:          f.OrderUID,
			TrackNumber:       f.TrackNumber,
			Entry:             f.Entry,
			Delivery:          entity.Delivery(f.Delivery),
			Payment:           entity.Payment(f.Payment),
			Items:             make([]entity.Item, len(f.Items)),
			Locale:            f.Locale,
			InternalSignature: f.InternalSignature,
			CustomerID:        f.CustomerID,
			DeliveryService:   f.DeliveryService,
			Shardkey:          f.Shardkey,
			SMID:              f.SMID,
			DateCreated:       f.DateCreated,
			OofShard:          f.OofShard,
		}

		order.Payment.Transaction = order.OrderUID

		for i, itemFake := range f.Items {
			order.Items[i] = entity.Item(itemFake)
		
			order.Items[i].TrackNumber = order.TrackNumber
		}
		order.Delivery.Phone = "+" + order.Delivery.Phone
		fakes = append(fakes, order)

	}

	return fakes
}
