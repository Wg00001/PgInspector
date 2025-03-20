package config

import "fmt"

/**
 * @description:
 * @author Wg
 * @date 2025/2/10
 */

type Id interface {
	GetIdentity() string
}

func (n Identity) GetIdentity() string {
	return n.Str()
}

func (n Identity) Str() string {
	return string(n)
}

func NewIdentity[T any](arg T) Identity {
	return Identity(fmt.Sprintf("%v", arg))
}
