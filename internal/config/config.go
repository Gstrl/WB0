package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	HTTPServer   `yaml:"http_server"`
	DBConnection `yaml:"db_connection"`
	NatsServer   `yaml:"nats_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DBConnection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type NatsServer struct {
	Url string `yaml:"url"`
}

// MustLoad В реальном проекте лучше читать из переменных окружения
func MustLoad() *Config {
	var cfg Config
	// Получаем текущую директорию
	dir, err := os.Getwd()
	fmt.Println(dir)
	if err != nil {
		fmt.Println("Ошибка получения текущего каталога:", err)

	}

	file, err := os.Open(filepath.Join(dir, "internal/config/local.yaml"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	// Читаем yaml
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Ошибка при чнении файла конфигурации")
	}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatal("Ошибка при парсинге файла конфигурации")
	}

	return &cfg
}
