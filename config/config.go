package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/utils/clipboard"
	"github.com/bndw/pick/utils/pswdgen"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	defaultConfigFileTmpl = "%s/.pick/config.toml"
)

type Config struct {
	Encryption crypto.Config
	Storage    backends.Config
	General    generalConfig
	Version    string
}

type generalConfig struct {
	Password pswdgen.Config
	// Warning: Deprecated. The PasswordLen field is required for backwards-compatibility :(
	PasswordLen int
	Clipboard   clipboard.Config
}

func Load(version string) (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	configFile := fmt.Sprintf(defaultConfigFileTmpl, home)
	config := Config{
		Storage:    backends.NewDefaultConfig(),
		Encryption: crypto.NewDefaultConfig(),
		General: generalConfig{
			Password:  pswdgen.NewDefaultConfig(),
			Clipboard: clipboard.NewDefaultConfig(),
		},
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
	// Warning: Deprecated. The PasswordLen field is required for backwards-compatibility :(
	if l := config.General.PasswordLen; l > 0 {
		config.General.Password.Length = l
	}

	return &config, nil
}
