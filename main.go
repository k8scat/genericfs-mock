package main

import (
	"flag"
	"log"

	"github.com/wanhuasong/genericfs/router"
	"github.com/wanhuasong/genericfs/utils"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	initFlags()

	if err := router.Run(); err != nil {
		panic(err)
	}
}

func initFlags() {
	flag.StringVar(&utils.Pubkey, "pubkey", "", "RSA pubkey path")
	flag.Parse()

	if utils.Pubkey == "" {
		log.Fatalln("pubkey not set")
	}
}
