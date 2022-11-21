package config

import (
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
}

func parseConfig(bytes []byte) (*Config, error) {
	var c Config
	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return nil, errors.Wrapf(err, "failed unmarshalling yaml")
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return nil, errors.Wrapf(err, "error validating config")
	}

	return &c, nil
}

func ReadConfig(cfgFile string) (*Config, error) {
	bytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading config file: %s", cfgFile)
	}

	cfg, err := parseConfig(bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed parsing config")
	}

	return cfg, nil
}
