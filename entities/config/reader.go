package config

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/4
 */

type Reader interface {
	NewReader(option map[string]string) (Reader, error)
	ReadConfig() error
	ReadInspector() error
	ReadAgent() error
	SaveIntoConfig()
	Watch() //todo
}
