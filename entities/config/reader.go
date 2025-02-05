package config

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/4
 */

type Reader interface {
	ReadFromSource() error
	SaveIntoConfig()
}

func InitConfig(reader Reader) {
	err := reader.ReadFromSource()
	if err != nil {
		panic(err)
	}
	reader.SaveIntoConfig()
}
