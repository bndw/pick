package commands

import (
	"errors"
	"fmt"

	pickErrors "github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	defaultNoteAction = "edit"
)

func init() {
	cmd := &cobra.Command{
		Use:   "note [name]",
		Short: "Create a note",
		Long:  "The note command is used to create, modify, list and remove secure notes.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Note, cmd, args)
		},
	}
	cmd.Flags().Bool("edit", false, "create or modify a note")
	cmd.Flags().Bool("ls", false, "list all notes")
	cmd.Flags().Bool("export", false, "export all notes")
	cmd.Flags().Bool("rm", false, "remove a note")
	rootCmd.AddCommand(cmd)
}

func Note(args []string, flags *pflag.FlagSet) error {
	name, err := parseNoteArgs(args)
	if err != nil {
		return err
	}

	action, err := parseNoteFlags(flags)
	if err != nil {
		return err
	}

	safe, err := newSafeLoader().Load()
	if err != nil {
		return err
	}

	switch action {
	case "edit":
		return safe.Notes.Edit(name)
	case "ls", "export":
		notes := safe.Notes.List()
		if len(notes) == 0 {
			return errors.New("No notes found")
		}
		if action == "ls" {
			for name, note := range notes {
				fmt.Printf("%s — last edited on %s\n", name, utils.FormatUnixTime(note.ModifiedOn))
			}
		} else {
			utils.PrettyPrint(notes)
		}
		return nil
	case "rm":
		return safe.Notes.Remove(name)
	}
	return nil
}

func parseNoteArgs(args []string) (name string, err error) {
	if len(args) > 1 {
		err = &pickErrors.InvalidCommandUsage{}
		return
	}
	if len(args) == 0 {
		return
	}
	name = args[0]

	return
}

func parseNoteFlags(flags *pflag.FlagSet) (string, error) {
	availableFlags := []string{
		"edit", "ls", "export", "rm",
	}
	for _, flagName := range availableFlags {
		action, err := flags.GetBool(flagName)
		if err != nil || action {
			return flagName, err
		}
	}

	return defaultNoteAction, nil
}
