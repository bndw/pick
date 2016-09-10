package commands

import (
	"github.com/bndw/pick/config"
	"github.com/spf13/cobra"
)

// Global pick config
var Config *config.Config

var rootCmd = &cobra.Command{
	Use:   "pick",
	Short: "pick is a minimal password manager",
	Long:  "pick is a minimal password manager",
}

func Execute(cfg *config.Config) error {
	Config = cfg

	return rootCmd.Execute()
}
