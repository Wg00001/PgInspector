package config

/**
 * @description:
 * @author Wg
 * @date 2025/2/10
 */

// todo:废弃此文件

type NameGetter interface {
	GetName() Name
}

func (n Name) GetName() Name {
	return n
}

func (n DBConfig) GetName() Name {
	return n.Name
}

func (n TaskConfig) GetName() Name {
	return n.Name
}

func GetName(n NameGetter) Name {
	return n.GetName()
}

func GetNameT(origin any) Name {
	switch t := origin.(type) {
	case string:
		return Name(t)
	case Name:
		return t
	case NameGetter:
		return t.GetName()
	}
	return ""
}

func (n Name) Str() string {
	return string(n)
}
