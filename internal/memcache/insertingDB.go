package memcache

import (
	"WB0/internal/order_struct"
	"database/sql"
)

func (c *Cache) InsertingDB(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM buy")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Итерация по результатам запроса и заполнение кеша
	for rows.Next() {
		var order order_struct.Order
		var id_order int
		err := rows.Scan(
			&id_order,
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSign,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard)
		if err != nil {
			panic(err)
		}

		var delivery order_struct.Delivery

		err = db.QueryRow("SELECT name_receiver, phone_number, zip_code, name_city, address, region, email FROM delivery WHERE delivery_id = $1", id_order).Scan(
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			panic(err)
		}
		order.Delivery = delivery

		var payment order_struct.Payment

		err = db.QueryRow("SELECT transaction_id, request_id, name_currency, name_provider, total_amount, payment_dt, name_bank, delivery_cost, goods_total, custom_fee FROM payment WHERE payment_id = $1", id_order).Scan(
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			panic(err)
		}

		order.Payment = payment

		rows_item, err := db.Query("SELECT chrt_id, track_number, price, rid, name_item, sale, size_item, total_price, nm_id, brand, status FROM item WHERE order_id = $1", id_order)
		if err != nil {
			panic(err)
		}
		defer rows_item.Close()

		var itemsOrder []order_struct.Item

		// Итерация по результатам запроса и заполнение кеша
		for rows_item.Next() {
			var item order_struct.Item
			err := rows_item.Scan(
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.RID,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status)
			if err != nil {
				panic(err)
			}
			itemsOrder = append(itemsOrder, item)
		}
		order.Items = itemsOrder
		c.Set(id_order, order, 0)
	}

	return nil
}
