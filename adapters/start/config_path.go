package start

import "github.com/wg00001/wgo-sdk/wg"

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/17
 */

var (
	pConfigType = "local_file"
	pFilePath   = wg.GetRelativePath("/app/config")
)

func SetConfigPath(filePath, configType string) {
	pConfigType = configType
	pFilePath = filePath
}
