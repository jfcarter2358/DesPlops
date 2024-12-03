package kustomize

import (
	"bytes"
	"desplops/sops"
	"fmt"
	"os"
	"os/exec"
)

func Template(overlay, kustomizePath string) (string, error) {
	cmd := exec.Command("kubectl", "kustomize", fmt.Sprintf("%s/%s", kustomizePath, overlay))
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outb.String(), nil
}

func WriteSecrets(config sops.SopsData, kustomizePath string) error {
	for overlay, values := range config.Secrets {
		for key, val := range values.(map[interface{}]interface{}) {
			if err := os.WriteFile(fmt.Sprintf("%s/%s/files/%v", kustomizePath, overlay, key), []byte(fmt.Sprintf("%v", val)), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
