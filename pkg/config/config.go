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
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DBConnection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func MustLoad() *Config {
	var cfg Config
	// Получаем текущую директорию
	dir, err := os.Getwd()
	fmt.Println(dir)
	if err != nil {
		fmt.Println("Error getting current directory:", err)

	}

	file, err := os.Open(filepath.Join(dir, "pkg/config/local.yaml"))
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
