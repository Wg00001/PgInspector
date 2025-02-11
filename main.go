package main

import (
	"PgInspector/adapters/config_reader"
	"PgInspector/entities/config"
	"PgInspector/usecase"
	"PgInspector/usecase/db"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

func main() {
	config.InitConfig(config_reader.NewReader("yaml", "app/config"))
	conn := db.Connect(usecase.GetDbConfig("example1"))
	if conn.Error() != nil {
		fmt.Println("\n", conn.Error())
		return
	}
	var res1, res2 string
	row := conn.QueryRow(usecase.GetInsp("5.5-2").SQL)
	err := row.Scan(&res1, &res2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("res: %+v %+v", res1, res2)
}
