package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
	"github.com/mitchellh/go-homedir"
)

const (
	defaultConfigFileTmpl = "%s/.pick/config.toml"
)

type Config struct {
	Encryption crypto.Config
	Storage    backends.Config
	Version    string
}

func Load(version string) (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	configFile := fmt.Sprintf(defaultConfigFileTmpl, home)
	config := Config{
		Encryption: crypto.NewDefaultConfig(),
	}
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			// TODO(): No config found, should we create one?
		} else {
			return nil, err
		}
	} else {
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			return nil, err
		}
	}

	config.Version = version

	return &config, nil
}
