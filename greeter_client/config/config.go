package config

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"os"
	"testgrpc/common"
)

type config struct {
	ServerUrl    string `toml:"server_url" name:"u"`
	ServerPort   uint   `toml:"serer_port" name:"p"`
	ClientName   string `toml:"client_name" name:"n"`
	DisableSSL   bool   `toml:"disable_ssl" name:"disablessl"`
	MutualAuth   bool   `toml:"mutual_auth" name:"mutualAuth"`
	SayHello     bool   `toml:"say_hello" name:"sayhello"`
	DownloadFile bool   `toml:"download_file" name:"downloadfile"`
	DownloadPath string `toml:"download_path" name:"downloadpath"`
	Cert         string `toml:"cert" name:"cert"`
	Key          string `toml:"key" name:"key"`
	CACert       string `toml:"ca_cert" name:"cacert"`
}

func LoadConfig(fileConfig string) (common.ConfMap, error) {
	var c config
	var Conf = make(common.ConfMap)
	configFile, err := os.Open(fileConfig)
	if err != nil {
		return Conf, errors.New(fmt.Sprintf("error loading config:%v", err))
	}
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return Conf, errors.New(fmt.Sprintf("error reading config:%v", err))
	}
	err = toml.Unmarshal(configBytes, &c)
	if err != nil {
		return Conf, errors.New(fmt.Sprintf("error unmarshalling config:%v", err))
	}
	common.MakeConfMap(&Conf, c)
	return Conf, nil
}
