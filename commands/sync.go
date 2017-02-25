package commands

import (
	"github.com/bndw/pick/backends"
	"github.com/bndw/pick/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "sync [other-pick-safe]",
		Short: "Sync current safe with another pick safe",
		Long:  "The sync command is used to sync the current pick safe with another pick safe.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Sync, cmd, args)
		},
	})
}

func Sync(args []string, flags *pflag.FlagSet) error {
	otherSafePath, err := parseSyncArgs(args)
	if err != nil {
		return err
	}

	safeLoader := newSafeLoader()

	safe, err := safeLoader.Load()
	if err != nil {
		return err
	}
	safeLoader.RememberPassword()

	otherStorageConfig := config.Storage
	// TODO(leon): :(
	otherStorageConfig.Settings["path"] = otherSafePath
	otherBackendClient, err := backends.NewDiskBackend(otherStorageConfig)
	if err != nil {
		return err
	}
	otherSafe, err := safeLoader.LoadWithBackendClient(otherBackendClient)
	if err != nil {
		return err
	}

	if err := safe.SyncWith(otherSafe); err != nil {
		return err
	}
	return nil
}

func parseSyncArgs(args []string) (otherSafePath string, err error) {
	if len(args) != 1 {
		err = &errors.InvalidCommandUsage{}
		return
	}

	otherSafePath = args[0]

	return
}
