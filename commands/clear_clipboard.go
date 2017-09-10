package commands

import (
	"strconv"
	"time"

	"github.com/bndw/pick/errors"
	"github.com/bndw/pick/utils/clipboard"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "clear-clipboard [after-seconds] [must-match]",
		Short: "Clear clipboard",
		Long:  "The clear-clipboard command is used to clear the clipboard.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(ClearClipboard, cmd, args)
		},
		Hidden: true,
	})
}

func ClearClipboard(args []string, flags *pflag.FlagSet) error {
	duration, match, err := parseClearClipboardArgs(args)
	if err != nil {
		return err
	}
	time.Sleep(duration)
	return clipboard.ClearIfMatch(match)
}

func parseClearClipboardArgs(args []string) (duration time.Duration, match string, err error) {
	if len(args) != 2 {
		err = errors.ErrInvalidCommandUsage
		return
	}

	var secs int64
	secs, err = strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return
	}
	duration = time.Duration(secs) * time.Second

	match = args[1]

	return
}
