package helm

import (
	"bytes"
	"desplops/sops"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

func Template(appName, chartPath, valuesPath string) (string, error) {
	cmd := exec.Command("helm", "template", strings.Replace(appName, "-", "", -1), chartPath, "-f", valuesPath)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return outb.String(), nil
}

func WriteSecrets(config sops.SopsData, chartPath, valuesPath string) error {
	valuesBytes, err := os.ReadFile(fmt.Sprintf("%s/values.yaml", chartPath))
	if err != nil {
		return err
	}
	var values map[string]interface{}
	if err := yaml.Unmarshal(valuesBytes, &values); err != nil {
		return err
	}
	values = recursiveWrite(config.Secrets, values)
	out, err := yaml.Marshal(values)
	return os.WriteFile(valuesPath, out, 0644)
}

func recursiveWrite(data map[string]interface{}, values map[string]interface{}) map[string]interface{} {
	for key, val := range data {
		if m, ok := val.(map[string]interface{}); ok {
			if v, ok := values[key]; ok {
				values[key] = recursiveWrite(m, v.(map[string]interface{}))
				continue
			}
		}
		values[key] = val
	}
	return values
}
