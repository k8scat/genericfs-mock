package utils

import (
	"fmt"
	"log"
	"testing"

	"github.com/wanhuasong/genericfs/config"
)

var (
	origin = "e=1625645423&t=1625641823&uuid=ButJmRBi"
	sig    string
)

func TestSign(t *testing.T) {
	Privkey = "/home/hsowan/workspace/genericfs/private-key.pem"
	log.Printf("Origin: %s", origin)

	var err error
	sig, err = Sign(origin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sig: %s\n", sig)
	fmt.Printf("sig len: %d\n", len(sig))
}

func TestVerifySig(t *testing.T) {
	origin = "e=1625469458&hash=Fk2wuo7kLZUYDWgJrY7FaGhHvWXJ&op=imageMogr2/auto-orient/thumbnail/x32&t=1625465858"
	sig = "5298797ba794b5ff6ad2a91b3011a67e127fca0c20a8a97f829729d8f272abdc599e8de6de404727864784a62f4493940a8d4f02cf6bc055fe26770974a75ad3e6e4c167c1c82db9526a430c2a0a0abf072fddf9ad8070c6285d2cdf84b42bc1e757e4999d67b76629192dbea4687ab42fc93f06b63c329dd35a28c84913c1a3"
	config.Cfg.PublicKey = "/home/hsowan/workspace/genericfs/public-key.pem"
	if err := VerifySig(origin, sig); err != nil {
		panic(err)
	}
}
