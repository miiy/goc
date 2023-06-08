package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func Load(filename string, config any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}
