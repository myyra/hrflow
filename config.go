package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/myyra/hrflow/hrflow"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func configPath() (string, error) {

	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("$HOME is not defined")
	}

	configFile := filepath.Join(dir, ".hrflow")

	return configFile, nil
}

func checkConfig(c *cli.Context) error {

	configFile, err := configPath()
	if err != nil {
		return errors.Wrap(err, "getting config path")
	}

	info, err := os.Stat(configFile)
	if os.IsNotExist(err) || info.IsDir() {
		return fmt.Errorf("can't find config file at %s", configFile)
	}

	return nil
}

func clientFromConfig() (*hrflow.Client, error) {

	configPath, err := configPath()
	if err != nil {
		return nil, errors.Wrap(err, "getting config path")
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "opening config file")
	}
	defer configFile.Close()

	var cfg config

	err = yaml.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "decoding config")
	}

	return hrflow.NewClient(cfg.Username, cfg.Password), nil
}
