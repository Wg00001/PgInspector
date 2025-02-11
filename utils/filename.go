package utils

import "strings"

func FileNameFormat(filename, configType string) string {
	if !strings.HasSuffix(filename, "."+configType) {
		filename += "." + configType
	}
	return filename
}
