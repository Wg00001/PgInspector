package main

import (
	"PgInspector/adapters/config_reader"
	"PgInspector/entities/config"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

func main() {
	config.InitConfig(config_reader.BuildReader("yaml", "app/config/config.yaml"))
}
