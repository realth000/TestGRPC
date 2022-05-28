package config

import "os"
import "github.com/pelletier/go-toml/v2"

type ClientConfig struct {
	ServerUrl    string `toml:"server_url"`
	ServerPort   string `toml:"server_port"`
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

func LoadConfigFile(filePath string, clientConfig *ClientConfig) error {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(configFile, clientConfig)
	return err
}
