package config

/**
 * @description:
 * @author Wg
 * @date 2025/2/10
 */

type Id interface {
	Identity() Id
}

func (n Identity) Identity() Id {
	return n
}

func (n Identity) Str() string {
	return string(n)
}
