package utils

import (
	"fmt"
	"testing"
)

var (
	origin = "e=1624250893&t=1624250893&hash=xxx"
	sig    string
)

func TestSignAndVerifySig(t *testing.T) {
	Privkey = "/home/hsowan/workspace/genericfs/private-key.pem"
	var err error
	sig, err = Sign(origin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sig: %s\n", sig)
	fmt.Printf("sig len: %d\n", len(sig))
}

func TestVerifySig(t *testing.T) {
	Pubkey = "/home/hsowan/workspace/genericfs/public-key.pem"
	if err := VerifySig(origin, sig); err != nil {
		panic(err)
	}
}
