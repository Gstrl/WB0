package consumerNats

import (
	"WB0/pkg/memcache"
	. "WB0/pkg/order_struct"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func RunConsumer(c *memcache.Cache) error {

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

		c.Set(c.AutoIncrement(), order, 0)

	})

	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()
	select {}
}
