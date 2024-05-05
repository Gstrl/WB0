package main

import (
	"WB0/pkg/HTTPServer"
	"WB0/pkg/config"
	"WB0/pkg/consumerNats"
	"WB0/pkg/db_connection"
	"WB0/pkg/memcache"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//1) Чтение конфига
	cfg := config.MustLoad()
	// 2) Подключение к Postgresql
	db, err := db_connection.DBConnect(cfg.DBConnection)
	if err != nil {
		log.Fatalf("Ошибка соединения с базой данных: %v", err)
	}

	//3)Инициализация cache
	cache := memcache.New(0, 0)
	//4)Записывем в кэш значения из базы данных
	err = cache.InsertingDB(db)
	if err != nil {
		log.Fatalf("Ошибка записи в кэш: %v", err)
	}
	//5)Подключение к NATS  серверу
	go func() {
		err := consumerNats.RunConsumer(db, cache)
		if err != nil {
			log.Fatal(err)
		}
	}()
	//6)Запуск сервера
	go func() {
		err := HTTPServer.RunServer(cache, cfg.Address)
		if err != nil {
			log.Fatal(err)
		}
	}()
	//7) Ожидание сигнала для завершения работы приложения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("Consumer завершает работу.")

}
