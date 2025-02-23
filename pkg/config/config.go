package config

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BurntSushi/toml"
)

var (
	mu  sync.RWMutex
	cfg *Config
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
	Id          string  `toml:"id"`
	ModelName   string  `toml:"model_name"`
	BaseUrl     string  `toml:"base_url"`
	ApiKey      string  `toml:"api_key"`
	Temperature float32 `toml:"temperature"`
	MaxTokens   int     `toml:"max_tokens"`
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
	cfgLocal := &Config{}
	_, err := toml.DecodeFile(cfgPath, cfgLocal)
	if err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()
	cfg = cfgLocal
	return cfg, nil
}

func WatchConfig(cfgPath string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)

	go func() {
		for {
			select {
			case <-sigChan:
				_, _ = LoadConfig(cfgPath)
			}
		}
	}()
}

func GetConfig() *Config {
	mu.RLock()
	defer mu.RUnlock()
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
