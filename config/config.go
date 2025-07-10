package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenAddr  string   `yaml:"listen_addr"`
	ReceivePath string   `yaml:"receive_path"`
	ForwardURLs []string `yaml:"forward_urls"`
}

var GlobalConfig *Config

func LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	GlobalConfig = &cfg
	return nil
}
