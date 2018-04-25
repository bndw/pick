package config

import (
	"fmt"

	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/crypto"
	"github.com/bndw/pick/utils/clipboard"
	"github.com/bndw/pick/utils/pswdgen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func Load(rootCmd *cobra.Command, version string) (*Config, error) {
	config := Config{
		Storage:    backends.NewDefaultConfig(),
		Encryption: crypto.NewDefaultConfig(),
		General: generalConfig{
			Password:  pswdgen.NewDefaultConfig(),
			Clipboard: clipboard.NewDefaultConfig(),
		},
		Version: version,
	}

	// Warning: Deprecated. The PasswordLen field is required for backwards-compatibility :(
	if l := config.General.PasswordLen; l > 0 {
		config.General.Password.Length = l
	}

	viper.AddConfigPath("$HOME/.pick")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	rootCmd.PersistentFlags().String("safe", "", "Overwrite path to safe file")
	viper.BindPFlag("storage.settings.path", rootCmd.PersistentFlags().Lookup("safe")) // file backend
	viper.BindPFlag("storage.settings.key", rootCmd.PersistentFlags().Lookup("safe"))  // s3 backend

	// Ugly, I know. See https://github.com/spf13/viper/issues/472
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := viper.ReadInConfig(); err != nil {
			// Return `nil` to avoid printing an unneccesary error message
			return nil
		}
		if err := viper.Unmarshal(&config); err != nil {
			return fmt.Errorf("failed to unmarshal into config: %v", err)
		}
		return nil
	}

	return &config, nil
}
