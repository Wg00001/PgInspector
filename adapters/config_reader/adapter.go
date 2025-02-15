package config_reader

import (
	"PgInspector/adapters/config_reader/local_yaml"
	"PgInspector/entities/config"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

func NewReader(configType, filepath string) config.Reader {
	switch configType {
	case "yaml":
		return &local_yaml.ConfigReaderYaml{FilePath: filepath}
	default:
		panic("reader type not exist: " + configType)
	}
	return nil
}
