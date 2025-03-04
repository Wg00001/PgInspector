package config

import (
	"PgInspector/adapters/config/yaml"
	"PgInspector/entities/config"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

const (
	Yaml = "yaml"
)

func NewReader(configType, filepath string) config.Reader {
	switch configType {
	case Yaml:
		return &yaml.ConfigReaderYaml{FilePath: filepath}
	default:
		fmt.Printf("reader type not exist: %s\n", configType)
	}
	return nil
}
