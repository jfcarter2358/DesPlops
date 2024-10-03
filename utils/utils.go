package utils

import (
	"desplops/sops"
	"fmt"
	"os"
	"strings"
)

func ReplaceConfig(config sops.SopsData, manifestPath string) error {
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	manifest := string(manifestBytes)

	for key, val := range config.Config {
		manifest = strings.ReplaceAll(manifest, fmt.Sprintf("$%s", key), fmt.Sprintf("%v", val))
	}

	return os.WriteFile(manifestPath, []byte(manifest), 0644)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
