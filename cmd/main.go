package main

import (
	"WB0/pkg/HTTPServer"
	"WB0/pkg/config"
	"WB0/pkg/consumerNats"
	"WB0/pkg/db_connection"
	"WB0/pkg/memcache"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//1) Чтение конфига
	cfg := config.MustLoad()
	// 1) Подключение к Postgresql
	db, err := db_connection.DBConnect(cfg.DBConnection)
	if err != nil {
		log.Fatalf("Ошибка соединения с базой данных: %v", err)
	}
	fmt.Println(db)

	//Инициализация cache
	cache := memcache.New(0, 0)
	//Подключение к NATS  серверу
	go func() {
		err := consumerNats.RunConsumer(cache)
		if err != nil {
			log.Fatal(err)
		}
	}()
	//Запуск сервера
	go func() {
		err := HTTPServer.RunServer(cache, cfg.Address)
		if err != nil {
			log.Fatal(err)
		}
	}()
	// Ожидание сигнала для завершения работы приложения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("Consumer завершает работу.")

}
