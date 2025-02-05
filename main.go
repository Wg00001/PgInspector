package main

import (
	"PgInspector/adapters/config_reader/local_yaml"
	"PgInspector/entities/config"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

func main() {
	config.InitConfig(&local_yaml.ConfigReaderYaml{FilePath: "app/config/config.yaml"})
}
