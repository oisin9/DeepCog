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

type Provider struct {
	Name    string `toml:"name"`
	BaseUrl string `toml:"base_url"`
	ApiKey  string `toml:"api_key"`
}

type BaseModel struct {
	Id        string `toml:"id"`
	ModelName string `toml:"model_name"`
	BaseUrl   string `toml:"base_url"`
	ApiKey    string `toml:"api_key"`
}

type Model struct {
	Id            string `toml:"id"`
	OwnedBy       string `toml:"owned_by"`
	ThinkingModel string `toml:"thinking_model"`
	GenerateModel string `toml:"generate_model"`
	ApiKey        string `toml:"api_key"`
}

type Config struct {
	Server     Server      `toml:"server"`
	Providers  []Provider  `toml:"providers"`
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

func (cfg *Config) GetModel(id string) *Model {
	for _, model := range cfg.Models {
		if model.Id == id {
			return &model
		}
	}
	return nil
}

func (cfg *Config) GetBaseModel(id string) *BaseModel {
	for _, baseModel := range cfg.BaseModels {
		if baseModel.Id == id {
			return &baseModel
		}
	}
	return nil
}
