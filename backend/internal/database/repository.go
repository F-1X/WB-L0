package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"wb/backend/internal/domain/entity"
	"wb/backend/internal/domain/repository"
)

var _ repository.OrderDB = &DB{}

type DB struct {
	conn PostgresDB
}

func NewPostgesRepository(conn PostgresDB) (repository.OrderDB, error) {
	return DB{conn: conn}, nil
}

// GetOrder - достает заказ из базы
func (db DB) GetOrder(ctx context.Context, id string) (entity.Order, error) {
	resultChan := make(chan struct {
		order entity.Order
		err   error
	})

	go func() {
		var order entity.Order
		query := `
		SELECT o.order_uid,
		       o.track_number,
		       o.entry,
		       o.locale,
		       o.internal_signature,
		       o.customer_id,
		       o.delivery_service,
		       o.shardkey,
		       o.sm_id,
		       o.date_created,
		       o.oof_shard,
		       d.name AS delivery_name,
		       d.phone AS delivery_phone,
		       d.zip AS delivery_zip,
		       d.city AS delivery_city,
		       d.address AS delivery_address,
		       d.region AS delivery_region,
		       d.email AS delivery_email,
		       p.transaction AS payment_transaction,
		       p.request_id AS payment_request_id,
		       p.currency AS payment_currency,
		       p.provider AS payment_provider,
		       p.amount AS payment_amount,
		       p.payment_dt AS payment_payment_dt,
		       p.bank AS payment_bank,
		       p.delivery_cost AS payment_delivery_cost,
		       p.goods_total AS payment_goods_total,
		       p.custom_fee AS payment_custom_fee,
		       COALESCE(jsonb_agg(jsonb_build_object(
		           'chrt_id', i.chrt_id,
		           'track_number', i.track_number,
		           'price', i.price,
		           'rid', i.rid,
		           'name', i.name,
		           'sale', i.sale,
		           'size', i.size,
		           'total_price', i.total_price,
		           'nm_id', i.nm_id,
		           'brand', i.brand,
		           'status', i.status
		       ) ORDER BY i.id), '[]'::jsonb) AS items
		FROM orders o
		JOIN delivery d ON o.order_uid = d.order_uid
		JOIN payment p ON o.order_uid = p.transaction
		JOIN items i ON o.track_number = i.track_number
		WHERE o.order_uid = $1
		GROUP BY o.order_uid, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		         p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
		         p.bank, p.delivery_cost, p.goods_total, p.custom_fee;
	`

		row := db.conn.QueryRow(ctx, query, id)

		var itemsJSON []byte
		err := row.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SMID,
			&order.DateCreated,
			&order.OofShard,
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDt,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
			&itemsJSON,
		)
		if err != nil {
			resultChan <- struct {
				order entity.Order
				err   error
			}{order, err}
			return
		}

		err = json.Unmarshal(itemsJSON, &order.Items)
		if err != nil {
			log.Println("error unmarshalling items JSON:", err)
			resultChan <- struct {
				order entity.Order
				err   error
			}{order, err}
			return
		}

		resultChan <- struct {
			order entity.Order
			err   error
		}{order, nil}
	}()

	select {
	case result := <-resultChan:
		return result.order, result.err
	case <-ctx.Done():
		return entity.Order{}, ctx.Err()
	}

}

// InsertOrder доабвляет новый заказ
func (db DB) InsertOrder(ctx context.Context, order entity.Order) error {
	resultChan := make(chan error)

	go func() {
		defer close(resultChan)
		var exists bool

		checkQuery := `SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid=$1)`
		err := db.conn.QueryRow(ctx, checkQuery, order.OrderUID).Scan(&exists)
		if err != nil {
			resultChan <- fmt.Errorf("failed to check order existence: %w", err)
		}

		if exists {
			resultChan <- ErrOrderExists
		}

		tx, err := db.conn.Begin(ctx)
		if err != nil {
			resultChan <- fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		orderQuery := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		_, err = tx.Exec(ctx, orderQuery,
			order.OrderUID,
			order.TrackNumber,
			order.Entry,
			order.Locale,
			order.InternalSignature,
			order.CustomerID,
			order.DeliveryService,
			order.Shardkey,
			order.SMID,
			order.DateCreated,
			order.OofShard,
		)
		if err != nil {
			resultChan <- fmt.Errorf("failed to insert into orders: %w", err)
		}

		deliveryQuery := `
			INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err = tx.Exec(ctx, deliveryQuery,
			order.OrderUID,
			order.Delivery.Name,
			order.Delivery.Phone,
			order.Delivery.Zip,
			order.Delivery.City,
			order.Delivery.Address,
			order.Delivery.Region,
			order.Delivery.Email,
		)
		if err != nil {
			resultChan <- fmt.Errorf("failed to insert into delivery: %w", err)
		}

		paymentQuery := `
			INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
		_, err = tx.Exec(ctx, paymentQuery,
			order.OrderUID,
			order.Payment.RequestID,
			order.Payment.Currency,
			order.Payment.Provider,
			order.Payment.Amount,
			order.Payment.PaymentDt,
			order.Payment.Bank,
			order.Payment.DeliveryCost,
			order.Payment.GoodsTotal,
			order.Payment.CustomFee,
		)
		if err != nil {
			resultChan <- fmt.Errorf("failed to insert into payment: %w", err)
		}

		stmtName := "insert_item"
		_, err = tx.Prepare(ctx, stmtName, `
			INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`)
		if err != nil {
			resultChan <- fmt.Errorf("failed to prepare statement: %w", err)
		}

		for _, item := range order.Items {
			_, err := tx.Exec(ctx, stmtName,
				item.ChrtID,
				order.TrackNumber,
				item.Price,
				item.Rid,
				item.Name,
				item.Sale,
				item.Size,
				item.TotalPrice,
				item.NMID,
				item.Brand,
				item.Status,
			)
			if err != nil {
				resultChan <- fmt.Errorf("failed to execute prepared statement: %w", err)
			}
		}

		if err = tx.Commit(ctx); err != nil {
			resultChan <- fmt.Errorf("failed to commit transaction: %w", err)
		}

	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (db DB) GetOrdersWithLimitByOrder(ctx context.Context, limit int, order string, direction string) ([]entity.Order, error) {

	type resultStruct struct {
		orders []entity.Order
		err   error
	}
	resultChan := make(chan resultStruct)

	go func() {
		if order == "" {
			order = "date_created"
		}

		if direction == "" {
			direction = "ASC"
		}

		query := fmt.Sprintf(`
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		ORDER BY %s %s
		LIMIT $1
	`, order, direction)

		rows, err := db.conn.Query(ctx, query, limit)
		if err != nil {
			log.Println("failed in query, err:", err)
			resultChan <- resultStruct{orders: nil, err:  fmt.Errorf("failed to query statement: %w", err)}
		}
		defer rows.Close()
		var orders []entity.Order
		for rows.Next() {
			var order entity.Order
			err := rows.Scan(
				&order.OrderUID,
				&order.TrackNumber,
				&order.Entry,
				&order.Locale,
				&order.InternalSignature,
				&order.CustomerID,
				&order.DeliveryService,
				&order.Shardkey,
				&order.SMID,
				&order.DateCreated,
				&order.OofShard,
			)
			if err != nil {
				log.Println("failed to scan order, err:", err)
				resultChan <- resultStruct{orders: nil, err:  fmt.Errorf("failed to scan order: %w", err)}
			}

			deliveryQuery := `
			SELECT name, phone, zip, city, address, region, email
			FROM delivery
			WHERE order_uid=$1
		`
			deliveryRow := db.conn.QueryRow(ctx, deliveryQuery, order.OrderUID)
			err = deliveryRow.Scan(
				&order.Delivery.Name,
				&order.Delivery.Phone,
				&order.Delivery.Zip,
				&order.Delivery.City,
				&order.Delivery.Address,
				&order.Delivery.Region,
				&order.Delivery.Email,
			)
			if err != nil {
				log.Println("failed to scan delivery details, err:", err)
				resultChan <- resultStruct{orders: nil, err: fmt.Errorf("failed to scan delivery details: %w", err)}
			}

			paymentQuery := `
			SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
			FROM payment
			WHERE order_uid=$1
		`
			paymentRow := db.conn.QueryRow(ctx, paymentQuery, order.OrderUID)
			err = paymentRow.Scan(
				&order.Payment.Transaction,
				&order.Payment.RequestID,
				&order.Payment.Currency,
				&order.Payment.Provider,
				&order.Payment.Amount,
				&order.Payment.PaymentDt,
				&order.Payment.Bank,
				&order.Payment.DeliveryCost,
				&order.Payment.GoodsTotal,
				&order.Payment.CustomFee,
			)
			if err != nil {
				resultChan <- resultStruct{orders: nil, err:  fmt.Errorf("failed to scan payment details: %w", err)}
			}

			itemsQuery := `
			SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
			FROM items
			WHERE order_uid=$1
		`
			itemRows, err := db.conn.Query(ctx, itemsQuery, order.OrderUID)
			if err != nil {
				resultChan <- resultStruct{orders: nil, err:  fmt.Errorf("failed to get items: %w", err)}
			}
			defer itemRows.Close()

			var items []entity.Item
			for itemRows.Next() {
				var item entity.Item
				err := itemRows.Scan(
					&item.ChrtID,
					&item.TrackNumber,
					&item.Price,
					&item.Rid,
					&item.Name,
					&item.Sale,
					&item.Size,
					&item.TotalPrice,
					&item.NMID,
					&item.Brand,
					&item.Status,
				)
				if err != nil {
					resultChan <- resultStruct{orders: nil, err:  fmt.Errorf("failed to scan item: %w", err)}
				}
				items = append(items, item)
			}
			if err := itemRows.Err(); err != nil {
				resultChan <- resultStruct{orders: nil, err: fmt.Errorf("itemRows error: %w", err)}
			}

			order.Items = items
			orders = append(orders, order)
		}

		if rows.Err() != nil {
			resultChan <- resultStruct{orders: nil, err: fmt.Errorf("rows error: %w", rows.Err())}
		}

		rows.Close()

	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case result := <-resultChan:
		return result.orders, result.err

	}
}
