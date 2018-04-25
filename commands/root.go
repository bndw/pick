package commands

import (
	c "github.com/bndw/pick/config"
	"github.com/spf13/cobra"
)

// Global pick config
var config *c.Config

var rootCmd = &cobra.Command{
	Use:   "pick",
	Short: "pick is a minimal password manager",
	Long:  "pick is a minimal password manager",
}

func RootCmd() *cobra.Command {
	return rootCmd
}

func Execute(cfg *c.Config) error {
	config = cfg

	return rootCmd.Execute()
}
