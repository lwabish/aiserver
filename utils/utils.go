package utils

import (
	"path/filepath"
	"strings"
)

func IsAllowedExtension(fileName string, allowedUploadExtensions []string) bool {
	for _, extension := range allowedUploadExtensions {
		if strings.ToLower(filepath.Ext(fileName)) == extension {
			return true
		}
	}
	return false
}
