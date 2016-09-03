package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdAdd = &cobra.Command{
	Use:   "add [name] [username] [password]",
	Short: "Add a credential",
	Long: `The add command is used to add a new credential.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(addCommand(args...))
	},
}

var cmdCat = &cobra.Command{
	Use:   "cat [name]",
	Short: "Cat a credential",
	Long: `The cat command is used to print a credential to stdout.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("USAGE: cat [name]")
			os.Exit(1)
		}

		os.Exit(catCommand(args...))
	},
}

var cmdCopy = &cobra.Command{
	Use:   "cp [name]",
	Short: "Copy a credential to the clipboard",
	Long: `The copy command is used to copy a credential's password
to the clipboard.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("USAGE: copy [name]")
			os.Exit(1)
		}

		os.Exit(copyCommand(args...))
	},
}

var cmdExport = &cobra.Command{
	Use:   "export",
	Short: "Export decrypted credentials in JSON format",
	Long: `The export command is used to export decrypted credentials in JSON
format.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(exportCommand(args...))
	},
}

var cmdList = &cobra.Command{
	Use:   "ls",
	Short: "List all credentials",
	Long: `The list command is used to list the saved credentials.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(listCommand(args...))
	},
}

var cmdRemove = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove a credential",
	Long: `The remove command is used to remove a saved credential.
            `,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(removeCommand(args...))
	},
}

var rootCmd = &cobra.Command{Use: "pick"}

func init() {
	rootCmd.AddCommand(cmdAdd, cmdCat, cmdCopy, cmdExport, cmdList, cmdRemove)
}

func main() {
	rootCmd.Execute()
}
