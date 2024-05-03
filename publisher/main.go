package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"os"
	"time"
)

type Order struct {
	OrderUID        string   `json:"order_uid"`
	TrackNumber     string   `json:"track_number"`
	Entry           string   `json:"entry"`
	Delivery        Delivery `json:"delivery"`
	Payment         Payment  `json:"payment"`
	Items           []Item   `json:"items"`
	Locale          string   `json:"locale"`
	InternalSign    string   `json:"internal_signature"`
	CustomerID      string   `json:"customer_id"`
	DeliveryService string   `json:"delivery_service"`
	ShardKey        string   `json:"shardkey"`
	SMID            int      `json:"sm_id"`
	DateCreated     string   `json:"date_created"`
	OOFShard        string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	endlessPublish(ctx, nc, js)
}

func endlessPublish(ctx context.Context, nc *nats.Conn, js jetstream.JetStream) {
	data, err := os.ReadFile("publisher/model.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла JSON: %v", err)
	}

	// Распарсивание данных JSON в структуру данных
	var order Order
	err = json.Unmarshal(data, &order)
	if err != nil {
		log.Fatalf("Ошибка парсинга JSON: %v", err)
	}

	// Использование данных из структуры
	log.Printf("Order UID: %s", order.OrderUID)
	log.Printf("Track Number: %s", order.TrackNumber)

	// Отправка сообщений
	jsonData, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Ошибка сериализации JSON: %v", err)
	}

	if nc.Status() == nats.CONNECTED {
		if _, err := js.Publish(ctx, "test.subject", jsonData); err != nil {
			fmt.Println("pub error: ", err)
		}
	}
}
