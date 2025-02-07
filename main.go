package main

import (
	"PgInspector/adapters/config_reader"
	"PgInspector/entities/config"
	"PgInspector/entities/db"
	"PgInspector/usecase"

	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

func main() {
	config.InitConfig(config_reader.BuildReader("yaml", "app/config/config.yaml"))
	conn, err := db.Connect(usecase.GetDbConfig("example1"))
	if err != nil {
		fmt.Println("\n", err)
		return
	}
	var res string
	row := conn.DB.QueryRow("SELECT email FROM users WHERE id = 1")
	err = row.Scan(&res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("res: %+v", res)
}
