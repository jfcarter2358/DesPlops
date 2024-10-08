package kapp

import (
	"bytes"
	"desplops/constants"
	"desplops/utils"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

func Deploy(kappName, manifestPath string) error {
	kappName = strings.ReplaceAll(kappName, "_", "-")
	cmd := exec.Command("kapp", "deploy", "-a", kappName, "--diff-changes", "-f", manifestPath, "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DryRun(kappName, manifestPath string) error {
	kappName = strings.ReplaceAll(kappName, "_", "-")
	cmd := exec.Command("kapp", "deploy", "-a", kappName, "--diff-changes", "--diff-run", "-f", manifestPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Backup(kappName, backupPath string) error {
	kappName = strings.ReplaceAll(kappName, "_", "-")
	cmd := exec.Command("kapp", "inspect", "-a", kappName, "--raw")
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	backup, err := cleanupBackup(outb.String())
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, []byte(backup), 0644)
}

func cleanupBackup(backupManifest string) (string, error) {
	r := bytes.NewReader([]byte(backupManifest))

	dec := yaml.NewDecoder(r)

	output := []string{}
	var doc map[string]interface{}
	for dec.Decode(&doc) == nil {
		if !utils.Contains(constants.API_FILTER, doc["kind"].(string)) {
			out, err := yaml.Marshal(doc)
			if err != nil {
				return "", err
			}
			output = append(output, string(out))
		}
	}

	return strings.Join(output, "\n---\n"), nil
}
