package local_yaml

import "PgInspector/entities/config"

/**
 * @description: inspector
 * @author Wg
 * @date 2025/1/19
 */

type ConfigReaderYaml struct {
	filepath string
}

var _ config.Reader = (*ConfigReaderYaml)(nil)

func (c ConfigReaderYaml) ReadFromSource() {
	panic("implement me")
}

func (c ConfigReaderYaml) SaveIntoConfig() {
	panic("implement me")
}
