package commands

import (
	"errors"
	"fmt"

	"github.com/bndw/pick/utils"
	"github.com/spf13/pflag"
)

func SafePass(args []string, flags *pflag.FlagSet) error {
	if !utils.Confirm("Do you really want to change the master password of your pick safe?", false) {
		return errors.New("Aborted as requested")
	}

	safeLoader := newSafeLoader(true)

	safe, err := safeLoader.Load()
	if err != nil {
		return err
	}
	safeLoader.RememberPassword()

	newPassword, err := readMasterPassConfirmed(true)
	if err != nil {
		return err
	}
	if err := safe.ChangePassword(newPassword); err != nil {
		return err
	}

	fmt.Println("Master password updated successfully")
	return nil
}
