package db

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/6
 */

type ConnEntity interface {
	Connect()
	Error() error
}
