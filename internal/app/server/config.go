package server

import (
	"flag"
	"fmt"
	"log"

	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/db/files"
	"github.com/BillyBones007/my-url-shortener/internal/db/maps"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
	"github.com/BillyBones007/my-url-shortener/internal/hasher/randchars"
	"github.com/caarlos0/env"
)

// Конфигурация сервера
type Config struct {
	// Адрес запуска HTTP сервера (env $SERVER_ADDRESS)
	ServerAddress string `env:"SERVER_ADDRESS"`
	// Базовый адресс результирующего сокращенного url (env $BASE_URL)
	BaseURL string `env:"BASE_URL"`
	// Путь к файлу хранилища (env $FILE_STORAGE_PATH)
	StoragePATH string `env:"FILE_STORAGE_PATH"`
	// Объект хранилища, реализующий интерфейс db.DBase (env $FILE_STORAGE_PATH)
	Storage db.DBase
	// Объект хэшера, реализующий интерфейс hasher.URLHasher
	Hash hasher.URLHasher
}

// Флаги командной строки
type FlagVars struct {
	bu string // base url
	sa string // server address
	sp string // storage path
}

// Конфигурирует сервер в зависимости от флагов командной строки и переменных окружения.
// В приоритете значения из переменных окружения.
func ServerConfigurator(cfg *Config, flagV *FlagVars) {
	fmt.Println("Вход в конфигуратор..")
	// Парсим переменные окружения в Config
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Задаем флаги командной строки
	flag.StringVar(&flagV.bu, "b", "", "Base ULR. If empty, it is replaced by the address of the running server.")
	flag.StringVar(&flagV.sa, "a", ":8080", "Server Address. By default: localhost:8080.")
	flag.StringVar(&flagV.sp, "f", "default", "Requires the path to the storage file, otherwise the default storage is in memory.")
	flag.Parse()
	// Далее проверяем, какие флаги и переменные окружения установлены
	// И в соответствии с результатами переопределяем поля структуры Config
	if cfg.BaseURL == "" {
		cfg.BaseURL = flagV.bu
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = flagV.sa
	}
	if cfg.StoragePATH == "" {
		if flagV.sp == "default" {
			cfg.Storage = maps.NewStorage()
			cfg.Hash = randchars.URLHash{}
			log.Println("INFO: $FILE_STORAGE_PATH, and flag -f is empty, map storage is used")
			paramConfServerInfo(cfg)
			return
		}
		var err error
		cfg.Storage, err = files.NewStorage(flagV.sp)
		if err != nil {
			log.Printf("ERROR: file storage error %s\n", err)
			log.Println("INFO: map storage is used")
			cfg.Storage = maps.NewStorage()
			cfg.Hash = randchars.URLHash{}
			paramConfServerInfo(cfg)
			return
		}
		cfg.StoragePATH = flagV.sp
		cfg.Hash = randchars.URLHash{}
		paramConfServerInfo(cfg)
		return
	}
	cfg.Storage, err = files.NewStorage(cfg.StoragePATH)
	if err != nil {
		log.Printf("ERROR: file storage error %s\n", err)
		log.Println("INFO: map storage is used")
		cfg.Storage = maps.NewStorage()
		cfg.Hash = randchars.URLHash{}
	}
	cfg.Hash = randchars.URLHash{}
	paramConfServerInfo(cfg)
}

// Выводит информацию о конфигурации сервера
func paramConfServerInfo(cfg *Config) {
	fmt.Println("INFO: Параметры конфигурирования сервера:")
	fmt.Printf("Адрес сервера: %s\n", cfg.ServerAddress)
	fmt.Printf("Базовый результирующий URL: %s\n", cfg.BaseURL)
	fmt.Printf("Путь до файла-хранилища: %s\n", cfg.StoragePATH)
	fmt.Printf("Указатель на объект хранилища: %v\n", cfg.Storage)
}
