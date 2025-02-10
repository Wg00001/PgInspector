package main

import (
	"PgInspector/adapters/config_reader"
	"PgInspector/entities/config"
	"PgInspector/usecase/db"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

func main() {
	config.InitConfig(config_reader.BuildReader("yaml", "app/config/config.yaml"))
	conn := db.Connect(config.Name("example1"))
	if conn.Error() != nil {
		fmt.Println("\n", conn.Error())
		return
	}
	var res string
	row := conn.QueryRow("SELECT email FROM users WHERE id = 1")
	err := row.Scan(&res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("res: %+v", res)
}
