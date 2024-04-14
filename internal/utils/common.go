package utils

import (
	"fmt"
	"os"
	"slices"
)

func ReadConfigFile(path string) string {
	// Valid envs: local, dev, uat, prod
	envs := []string{"local", "dev", "uat", "prod"}

	env := os.Getenv("ENV")
	if env == "" || !slices.Contains(envs, env) {
		fmt.Println("No valid env provided, defaulting to local")
		env = "local"
	}

	filePath := fmt.Sprintf("%s/config_%s.yaml", path, env)
	return filePath
}
