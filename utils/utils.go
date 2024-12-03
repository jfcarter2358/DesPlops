package utils

import (
	"desplops/sops"
	"fmt"
	"strings"
)

func ReplaceConfig(config sops.SopsData, manifest string) string {
	for key, val := range config.Config {
		manifest = strings.ReplaceAll(manifest, fmt.Sprintf("$%s", key), fmt.Sprintf("%v", val))
	}

	return manifest
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
