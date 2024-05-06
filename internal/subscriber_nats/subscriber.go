package subscriber_nats

import (
	"WB0/internal/config"
	"WB0/internal/memcache"
	. "WB0/internal/order_struct"
	"database/sql"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

func RunSubscriber(db *sql.DB, c *memcache.Cache, cfg *config.Config) error {

	nc, err := nats.Connect(cfg.Url)
	if err != nil {
		return err
	}
	log.Println("Установка соединения с сервером NATS Streaming - успешно")
	defer nc.Close()

	var order Order
	subscription, err := nc.Subscribe("test.subject", func(msg *nats.Msg) {
		log.Printf("Получено сообщение: %s", string(msg.Data))

		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatalf("Ошибка парсинга JSON: %v", err)
		}

		//Добавление заказа в базу данных
		err := InsertOrder(db, order)
		if err != nil {
			log.Fatalf("Ошибка записи заказа в базу данных: %v", err)
			return
		} else { //Добавление заказа в кэш
			c.Set(c.AutoIncrement(), order, 0)
		}

	})

	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()
	select {}
}
