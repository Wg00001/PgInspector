package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// 从本地配置文件中加载并初始化配置

//todo:fix
func LoadConfigYaml(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	return &config, nil
}

func LoadConfigTable() {

}

func LoadConfigEtcd() {

}

func LoadConfigViper() {

}
