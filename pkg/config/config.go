package config

import (
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	once sync.Once
	cfg  *Config
)

type Server struct {
	Port string `toml:"port"`
}

type BaseModel struct {
	Name      string `toml:"name"`
	ModelName string `toml:"model_name"`
	BaseUrl   string `toml:"base_url"`
	ApiKey    string `toml:"api_key"`
}

type Model struct {
	Name          string `toml:"name"`
	ThinkingModel string `toml:"thinking_model"`
	GenerateModel string `toml:"generate_model"`
	ApiKey        string `toml:"api_key"`
}

type Config struct {
	Server     Server      `toml:"server"`
	BaseModels []BaseModel `toml:"base_models"`
	Models     []Model     `toml:"models"`
}

func LoadConfig(cfgPath string) (*Config, error) {
	loadConfigErr := error(nil)
	once.Do(func() {
		cfgLocal := &Config{}
		_, err := toml.DecodeFile(cfgPath, cfgLocal)
		if err != nil {
			loadConfigErr = err
			return
		}
		cfg = cfgLocal
	})
	return cfg, loadConfigErr
}

func GetConfig() *Config {
	return cfg
}
