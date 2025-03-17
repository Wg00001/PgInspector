package config

/**
 * @description:
 * @author Wg
 * @date 2025/2/10
 */

// todo:废弃此文件

type Identity interface {
	Identity() Name
}

func (n Name) Identity() Name {
	return n
}

func (n Name) Str() string {
	return string(n)
}
