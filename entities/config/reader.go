package config

import "fmt"

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/4
 */

type Reader interface {
	ReadConfig() error
	ReadInspector() error
	SaveIntoConfig()
	FormatFilename()
}

func InitConfig(reader Reader) {
	reader.FormatFilename()
	fmt.Println(reader)
	err := reader.ReadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", reader)
	err = reader.ReadInspector()
	if err != nil {
		panic(err)
	}
	reader.SaveIntoConfig()
}
