package config

import (
	"os"
)

func SetupMock() {
	basePath := os.Getenv("PROJECT_ROOT")
	if basePath == "" {
		panic(`environment variable "PROJECT_ROOT" undefined`)
	}
	getSearchPath = func() []string {
		searchPaths := []string{
			"/configs/test.local.yaml",
			"/configs/test.yaml",
		}
		for i := range searchPaths {
			searchPaths[i] = basePath + searchPaths[i]
		}
		return searchPaths
	}
}
