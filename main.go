package main

import (
	"log"

	"github.com/Skudarnov-Alexander/URLshortener/internal/api"
	"github.com/Skudarnov-Alexander/URLshortener/internal/logger"
	"github.com/Skudarnov-Alexander/URLshortener/internal/url/db"
)

func main() {
	Config, err := api.NewConf("./internal/api/config")
	if err != nil {
		log.Print("config reading error")
		return
	}
	log.Print("Config reading - DONE")

	Logger, err := logger.New(Config.LogPath)
	if err != nil {
		log.Print("Logger creation error" + err.Error()) // куда писать ошибку в создании логгера?)
		return
	}
	log.Print("Logger creation - DONE")

	InitDataBase := db.New()

	R := api.New(Logger, InitDataBase, Config)

	R.Start()

	// TODO грейсфул шатдаун

}
