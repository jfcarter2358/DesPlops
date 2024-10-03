package main

import (
	"desplops/helm"
	"desplops/kapp"
	"desplops/kustomize"
	"desplops/sops"
	"desplops/utils"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jfcarter2358/go-logger"
)

func main() {
	// setup CLI

	backupPath := flag.String("backup", "./backup-manifest.yaml", "Path to output backup manifest on deploy. Defaults to './backup-manifest.yaml'")
	configPath := flag.String("config", "config.yaml", "Path to sops config file. Defaults to 'config.yaml'")
	mode := flag.String("mode", "", "Path to write values file for Helm deployment. Valid values are 'kustomize' and 'helm'. Required")
	outputPath := flag.String("output", "./rendered-manifest.yaml", "Path to output rendered manifest. Defaults to './rendered-manifest.yaml'")
	rollback := flag.Bool("rollback", false, "Should rollback on deploy failure. Defaults to false")
	templatePath := flag.String("template", "", "Path to either the Kustomize overlay (e.g. ${PWD}/kustomize/overlays/foo) or the Helm chart directory containing the Chart.yaml file (e.g. ${PWD}/helm/foo). Required")
	valuesPath := flag.String("values", "./values.yaml", "Path to write values file for Helm deployment. Defaults to './values.yaml'")
	yes := flag.Bool("yes", false, "Should a non-dry run deploy be performed. Defaults to false")
	logLevel := flag.String("loglevel", logger.LOG_LEVEL_WARN, "Log level to use. Valid values are 'NONE', 'FATAL', 'SUCCESS', 'ERROR', 'WARN', 'INFO', 'DEBUG', and 'TRACE'. Defaults to 'WARN'")

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
	logger.Debugf("", "backup path        | %s", *backupPath)
	logger.Debugf("", "config path        | %s", *configPath)
	logger.Debugf("", "mode               | %s", *mode)
	logger.Debugf("", "output path        | %s", *outputPath)
	logger.Debugf("", "rollback           | %v", *rollback)
	logger.Debugf("", "template path      | %s", *templatePath)
	logger.Debugf("", "values path        | %s", *valuesPath)
	logger.Debugf("", "yes                | %v", *yes)
	logger.Tracef("", "app name           | %s", appName)
	logger.Tracef("", "base template path | %s", baseTemplatePath)

	config, err := sops.GetConfig(*configPath)
	if err != nil {
		panic(err)
	}

	// render manifest based on mode
	logger.Infof("", "Rendering with mode '%s'...", *mode)
	switch *mode {
	case "kustomize":
		logger.Tracef("", "Dropped into kustomize case")
		if err := kustomize.WriteSecrets(config, baseTemplatePath); err != nil {
			logger.Fatalf("", "Error on kustomize secrets write: %s", err.Error())
		}
		if err := kustomize.Template(appName, baseTemplatePath, *outputPath); err != nil {
			logger.Fatalf("", "Error on kustomize rendering: %s", err.Error())
		}
		if err := utils.ReplaceConfig(config, *outputPath); err != nil {
			logger.Fatalf("", "Error on manifest value replacement: %s", err.Error())
		}
	case "helm":
		logger.Tracef("", "Dropped into helm case")
		if err := helm.WriteSecrets(config, *templatePath, *valuesPath); err != nil {
			logger.Fatalf("", "Error on helm secrets write: %s", err.Error())
		}
		if err := helm.Template(appName, *templatePath, *valuesPath, *outputPath); err != nil {
			logger.Fatalf("", "Error on kustomize rendering: %s", err.Error())
		}
		if err := utils.ReplaceConfig(config, *outputPath); err != nil {
			logger.Fatalf("", "Error on manifest value replacement: %s", err.Error())
		}
	default:
		fmt.Printf("Invalid mode '%s', valid modes are 'helm' and 'kustomize'", *mode)
		flag.Usage()
		os.Exit(1)
	}
	logger.Success("", "Finished rendering!")

	// do the deploy
	if *yes {
		logger.Tracef("", "Dropped into yes=true case")
		logger.Infof("", "Backing up current deployment...")
		if err := kapp.Backup(appName, *backupPath); err != nil {
			logger.Fatalf("", "Error on kapp backup: %s", err.Error())
		}
		logger.Successf("", "Finished backup!")
		logger.Infof("", "Doing deploy...")
		if err := kapp.Deploy(appName, *outputPath); err != nil {
			if *rollback {
				logger.Tracef("", "Dropped into rollback=true case")
				logger.Errorf("", "Error on kapp deploy: %s", err.Error())
				logger.Infof("", "Doing rollback...")
				if err := kapp.Deploy(appName, *backupPath); err != nil {
					logger.Fatalf("", "Error on kapp rollback: %s", err.Error())
				}
				logger.Successf("", "Finished rollback!")
				os.Exit(0)
			} else {
				logger.Tracef("", "Dropped into rollback=false case")
				logger.Fatalf("", "Error on kapp deploy: %s", err.Error())
			}
		}
		logger.Successf("", "Finished deploy!")
	} else {
		logger.Tracef("", "Dropped into yes=false case")
		logger.Infof("", "Doing dry-run...")
		if err := kapp.DryRun(appName, *outputPath); err != nil {
			logger.Fatalf("", "Error on kapp dry-run: %s", err.Error())
		}
		logger.Successf("", "Finished dry-run!")
	}
}
