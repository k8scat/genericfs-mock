package main

import (
	"flag"
	"log"

	"github.com/wanhuasong/genericfs/config"
	"github.com/wanhuasong/genericfs/models"
	"github.com/wanhuasong/genericfs/router"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	initFlags()
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := models.InitDB(); err != nil {
		panic(err)
	}
	if err := router.Run(); err != nil {
		panic(err)
	}
}

func initFlags() {
	flag.StringVar(&config.CfgFile, "c", "./config.json", "Config file")
	flag.Parse()
}
