package config

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"os"
)

type Config struct {
	ServerUrl    string `toml:"server_url"`
	ServerPort   uint   `toml:"serer_port"`
	ClientName   string `toml:"client_name"`
	DisableSSL   bool   `toml:"disable_ssl"`
	MutualAuth   bool   `toml:"mutual_auth"`
	SayHello     bool   `toml:"say_hello"`
	DownloadFile bool   `toml:"download_file"`
	DownloadPath string `toml:"download_path"`
	Cert         string `toml:"cert"`
	Key          string `toml:"key"`
	CACert       string `toml:"ca_cert"`
}

func LoadConfig(fileConfig string) (Config, error) {
	var config Config
	configFile, err := os.Open(fileConfig)
	if err != nil {
		return Config{}, errors.New(fmt.Sprintf("error loading config:%v", err))
	}
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return Config{}, errors.New(fmt.Sprintf("error reading config:%v", err))
	}
	err = toml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, errors.New(fmt.Sprintf("error unmarshalling config:%v", err))
	}
	return config, nil
}
