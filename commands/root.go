package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "pick",
	Short: "pick is a minimal password manager",
	Long:  "pick is a minimal password manager",
}
