package sops

import (
	"bytes"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type SopsData struct {
	Secrets map[string]interface{} `yaml:"secrets" json:"secrets"`
	Config  map[string]interface{} `yaml:"config" json:"config"`
}

func GetConfig(configPath string) (SopsData, error) {
	cmd := exec.Command("sops", "decrypt", configPath)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return SopsData{}, err
	}

	var outData SopsData
	if err := yaml.Unmarshal(outb.Bytes(), &outData); err != nil {
		return SopsData{}, err
	}

	return outData, nil
}
