package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/wanhuasong/genericfs/config"
)

var (
	Privkey string
)

func VerifySig(origin, sig string) error {
	b, err := ioutil.ReadFile(config.Cfg.PublicKey)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(b)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hashed := SHA1([]byte(origin))
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hashed, []byte(sigBytes))
}

func Sign(s string) (string, error) {
	key, err := ioutil.ReadFile(Privkey)
	if err != nil {
		return "", err
	}
	r, err := encryptSHA1WithRSA(key, []byte(s))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(r), nil
}

func encryptSHA1WithRSA(key, data []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("no PEM data is found")
	}

	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	hashed := SHA1(data)
	return rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA1, hashed)
}

func SHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func SortParams(params map[string]string) string {
	keys := make([]string, 0)
	for key := range params {
		keys = append(keys, key)
	}
	var s string
	sort.Strings(keys)
	for _, key := range keys {
		if s != "" {
			s += "&"
		}
		s += fmt.Sprintf("%s=%s", key, params[key])
	}
	return s
}
