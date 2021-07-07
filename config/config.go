package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	Cfg     *Config
	CfgFile string
)

type Config struct {
	DSN       string `json:"dsn"`
	PublicKey string `json:"public_key"`
	StoreDir  string `json:"store_dir"`
}

func InitConfig() error {
	b, err := ioutil.ReadFile(CfgFile)
	if err != nil {
		return err
	}
	Cfg = new(Config)

	err = json.Unmarshal(b, &Cfg)
	if err != nil {
		return err
	}

	if Cfg.PublicKey == "" {
		return errors.New("public_key cannot be empty")
	}

	err = initStoreDir(Cfg.StoreDir)
	return err
}

func initStoreDir(p string) error {
	f, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(p, 644)
		}
		return err
	}
	if !f.IsDir() {
		err = os.MkdirAll(p, 644)
	}
	return err
}
