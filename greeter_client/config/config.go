package config

import "os"
import "github.com/pelletier/go-toml/v2"

type ClientConfig struct {
	ServerUrl        string `toml:"server_url"`
	ServerPort       uint   `toml:"server_port"`
	ClientName       string `toml:"client_name"`
	SSL              bool   `toml:"ssl"`
	SSLCert          string `toml:"ssl_cert"`
	SSLKey           string `toml:"ssl_key"`
	SSLCACert        string `toml:"ssl_ca_cert"`
	MutualAuth       bool   `toml:"mutual_auth"`
	SayHello         bool   `toml:"say_hello"`
	DownloadFile     bool   `toml:"download_file"`
	DownloadFilePath string `toml:"download_file_path"`
}

func LoadConfigFile(filePath string, clientConfig *ClientConfig) error {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(configFile, clientConfig)
	return err
}
