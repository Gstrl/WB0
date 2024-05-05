package consumerNats

import (
	"WB0/pkg/memcache"
	. "WB0/pkg/order_struct"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func RunConsumer(db *sql.DB, c *memcache.Cache) error {

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		return err
	}
	fmt.Println("Установка соединения с сервером NATS Streaming - успешно")
	defer nc.Close()

	var order Order
	subscription, err := nc.Subscribe("test.subject", func(msg *nats.Msg) {
		log.Printf("Получено сообщение: %s", string(msg.Data))

		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatalf("Ошибка парсинга JSON: %v", err)
		}

		//Добавление заказа в кэш
		c.Set(c.AutoIncrement(), order, 0)
		//Добавление заказа в базу данных
		err := InsertOrder(db, order)
		if err != nil {
			log.Fatalf("Ошибка записи заказа в базу данных: %v", err)
			return
		}

	})

	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()
	select {}
}

func InsertOrder(db *sql.DB, order Order) error {

	fmt.Println("вставляет данные о заказе в базу данных")

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Откатить транзакцию при любой ошибке

	// Вставка данных в таблицу buy
	_, err = tx.Exec(`
		INSERT INTO buy (order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSign, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OOFShard)
	if err != nil {
		return err
	}

	// Вставка данных в таблицу delivery
	_, err = tx.Exec(`
		INSERT INTO delivery (name_receiver, phone_number, zip_code, name_city, address, region, email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	// Вставка данных в таблицу payment
	_, err = tx.Exec(`
		INSERT INTO payment (transaction_id, request_id, name_currency, name_provider, 
		total_amount, payment_dt, name_bank, delivery_cost, goods_total, custom_fee) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}
	// Получаем идентификатор заказа
	var orderID int

	err = tx.QueryRow("SELECT MAX(buy_id) FROM buy WHERE order_uid = $1", order.OrderUID).Scan(&orderID)
	if err != nil {
		return err
	}
	fmt.Println(orderID)
	// Вставка данных в таблицу items
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO item (order_id, chrt_id, track_number, price, rid, name_item, sale, size_item, 
			total_price, nm_id, brand, status) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			orderID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
