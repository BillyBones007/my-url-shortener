package main

import (
	"flag"
	"log"

	"github.com/BillyBones007/my-url-shortener/internal/app/server"
	"github.com/BillyBones007/my-url-shortener/internal/db/files"
	"github.com/BillyBones007/my-url-shortener/internal/db/maps"
	"github.com/BillyBones007/my-url-shortener/internal/hasher/randchars"
	"github.com/caarlos0/env/v6"
)

type flagVars struct {
	bu string // base url
	sa string // server address
	sp string // storage path
}

var cfg server.Config
var fV flagVars

func init() {
	// Парсим переменные окружения в config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Задаем флаги командной строки
	flag.StringVar(&fV.bu, "b", "", "Base ULR. If empty, it is replaced by the address of the running server.")
	flag.StringVar(&fV.sa, "a", ":8080", "Server Address. By default: localhost:8080.")
	flag.StringVar(&fV.sp, "f", "default", "Requires the path to the storage file, otherwise the default storage is in memory.")
}

func main() {
	flag.Parse()
	serverConfigurator(&cfg, &fV)
	server := server.NewServer(&cfg)
	log.Fatal(server.ListenAndServe())
}

// Конфигурирует сервер в зависимости от флагов командной строки и переменных окружения.
// В приоритете значения из переменных окружения.
func serverConfigurator(cfg *server.Config, flagV *flagVars) {
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
			return
		}
		var err error
		cfg.Storage, err = files.NewStorage(flagV.sp)
		if err != nil {
			log.Printf("ERROR: file storage error %s\n", err)
			log.Println("INFO: map storage is used")
			cfg.Storage = maps.NewStorage()
			cfg.Hash = randchars.URLHash{}
			return
		}

	}
	cfg.Hash = randchars.URLHash{}
}
