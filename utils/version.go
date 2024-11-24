package utils

import "os"

// GetVersion read version from `version` file
func GetVersion() string {
	// theres a file named `version` in the root of the project (sibling to `main.go`)
	// which contains the version of the project
	// this function reads the content of that file and returns it
	verFilePath := "version"
	ver, err := os.ReadFile(verFilePath)
	if err != nil {
		return "unknown"
	}
	return string(ver)
}
