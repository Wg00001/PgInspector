package config

import (
	"log"
)

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
	err := reader.ReadConfig()
	if err != nil {
		panic(err)
	}
	err = reader.ReadInspector()
	if err != nil {
		panic(err)
	}
	reader.SaveIntoConfig()
	log.Println("config initiated...")
}
