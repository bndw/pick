package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize pick",
		Long:  "The init command is used to initialize pick.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Init, cmd, args)
		},
	})
}

func Init(args []string, flags *pflag.FlagSet) error {
	return initSafe()
}
