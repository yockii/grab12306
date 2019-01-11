package config

import (
	"flag"
	"sync"
)

// Config 结构
type Config struct {
}

var instance *Config
var once sync.Once
var configFile = flag.String("configFile", "config.toml", "General configuration file")

// GetInstance 获取config唯一实例
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}
