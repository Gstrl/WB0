package main

import (
	"WB0/pkg/config"
	"WB0/pkg/consumerNats"
	"WB0/pkg/memcache"
	"WB0/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Инициализация конфига
	cfg := config.MustLoad()
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
		err := server.RunServer(cache, cfg.Address)
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
