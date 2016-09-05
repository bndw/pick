package main

import (
	"fmt"

	"github.com/bndw/pick/commands"
	"github.com/spf13/cobra"
)

const Version = "v0.2.2"

func init() {
	commands.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of pick",
		Long:  `The version command prints the version of pick`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("pick %s\n", Version)
		},
	})
}
