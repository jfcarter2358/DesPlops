package main

import (
	"desplops/helm"
	"desplops/kustomize"
	"desplops/sops"
	"desplops/utils"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jfcarter2358/go-logger"
)

type configOverrides []string

func (c *configOverrides) String() string {
	return fmt.Sprintf("%s", *c)
}

func (c *configOverrides) Set(value string) error {
	parts := strings.Split(value, "=")

	if len(parts) != 2 {
		return fmt.Errorf("config override value must be of format '<key>=<val>'")
	}
	*c = append(*c, value)

	return nil
}

func main() {
	// setup CLI
	var c configOverrides

	configPath := flag.String("config", "config.yaml", "Path to sops config file. Defaults to 'config.yaml'")
	mode := flag.String("mode", "", "Path to write values file for Helm deployment. Valid values are 'kustomize' and 'helm'. Required")
	flag.Var(&c, "override", "Config values to override, takes format of '<key>=<val>")
	templatePath := flag.String("template", "", "Path to either the Kustomize overlay (e.g. ${PWD}/kustomize/overlays/foo) or the Helm chart directory containing the Chart.yaml file (e.g. ${PWD}/helm/foo). Required")
	valuesPath := flag.String("values", "./values.yaml", "Path to write values file for Helm deployment. Defaults to './values.yaml'")
	logLevel := flag.String("loglevel", logger.LOG_LEVEL_NONE, "Log level to use. Valid values are 'NONE', 'FATAL', 'SUCCESS', 'ERROR', 'WARN', 'INFO', 'DEBUG', and 'TRACE'. Defaults to 'NONE'")

	flag.Parse()

	logger.SetLevel(*logLevel)

	splitTemplatePath := strings.Split(*templatePath, "/")
	appName := splitTemplatePath[len(splitTemplatePath)-1]
	baseTemplatePath := strings.Join(splitTemplatePath[:len(splitTemplatePath)-1], "/")

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	logger.Debugf("", "Dumping options")
	logger.Debugf("", "------------------------------------------------")
	logger.Debugf("", "config path        | %s", *configPath)
	logger.Debugf("", "mode               | %s", *mode)
	logger.Debugf("", "")
	logger.Debugf("", "template path      | %s", *templatePath)
	logger.Debugf("", "values path        | %s", *valuesPath)
	logger.Tracef("", "base template path | %s", baseTemplatePath)

	config, err := sops.GetConfig(*configPath)
	if err != nil {
		panic(err)
	}

	for _, override := range c {
		parts := strings.Split(override, "=")
		config.Config[parts[0]] = parts[1]
	}

	// render manifest based on mode
	logger.Infof("", "Rendering with mode '%s'...", *mode)
	switch *mode {
	case "kustomize":
		logger.Tracef("", "Dropped into kustomize case")
		if err := kustomize.WriteSecrets(config, baseTemplatePath); err != nil {
			logger.Fatalf("", "Error on kustomize secrets write: %s", err.Error())
		}
		manifest, err := kustomize.Template(appName, baseTemplatePath)
		if err != nil {
			logger.Fatalf("", "Error on kustomize rendering: %s", err.Error())
		}
		manifest = utils.ReplaceConfig(config, manifest)
		fmt.Print(manifest)
	case "helm":
		logger.Tracef("", "Dropped into helm case")
		if err := helm.WriteSecrets(config, *templatePath, *valuesPath); err != nil {
			logger.Fatalf("", "Error on helm secrets write: %s", err.Error())
		}
		manifest, err := helm.Template(appName, *templatePath, *valuesPath)
		if err != nil {
			logger.Fatalf("", "Error on kustomize rendering: %s", err.Error())
		}
		manifest = utils.ReplaceConfig(config, manifest)
		fmt.Print(manifest)
	default:
		fmt.Printf("Invalid mode '%s', valid modes are 'helm' and 'kustomize'", *mode)
		flag.Usage()
		os.Exit(1)
	}
	logger.Success("", "Finished rendering!")
}
