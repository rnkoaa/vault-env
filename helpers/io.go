package helpers

import (
	"os"
	"strings"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsYamlFile - indicates if the file is a yaml file or not
func IsYamlFile(path string) bool {
	return strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")
}

// IsJSONFile - tells if a file is a json file
func IsJSONFile(path string) bool {
	return strings.HasSuffix(path, ".json")
}

// GetFileType - file type not found.
func GetFileType(path string) (string, error) {
	if !Exists(path) {
		return "", ErrorFileNotFound
	}
	if IsYamlFile(path) {
		return "yaml", nil
	}
	if IsJSONFile(path) {
		return "json", nil
	}

	return "", ErrorFileTypeUndetermined
}
