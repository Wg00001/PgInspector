package config

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/4
 */

type Reader interface {
	ReadFromSource()
	SaveIntoConfig()
}
