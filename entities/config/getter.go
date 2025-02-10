package config

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/10
 */

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
	return n.TaskName
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
