package commands

import (
	"fmt"

	"github.com/bndw/pick/utils/pswdgen"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	cmd := &cobra.Command{
		Use:   "pass",
		Short: "Generate a password without storing it",
		Long:  "The pass command is used to generate a password without storing it.",
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(Pass, cmd, args)
		},
	}
	cmd.Flags().Int("length", pswdgen.DefaultPasswordLength, "the length of the generated password")
	cmd.Flags().String("strength", "full", "the strength of the generated password")
	cmd.Flags().Int("num", 1, "the number of passwords to generate")
	rootCmd.AddCommand(cmd)
}

func Pass(args []string, flags *pflag.FlagSet) error {
	length, strength, num, err := parsePassFlags(flags)
	if err != nil {
		return err
	}
	config.General.Password.Length = length
	config.General.Password.Strength = pswdgen.StrengthByString(strength)
	config.General.Password.Mode = pswdgen.PasswordModeInteractive
	if num > 1 {
		config.General.Password.Mode = pswdgen.PasswordModeNonInteractive
	}

	for i := 0; i < num; i++ {
		password, err := pswdgen.Generate(config.General.Password)
		if err != nil {
			return err
		}
		fmt.Println(password)
	}
	return nil
}

func parsePassFlags(flags *pflag.FlagSet) (length int, strength string, num int, err error) {
	if length, err = flags.GetInt("length"); err != nil {
		return
	}
	if strength, err = flags.GetString("strength"); err != nil {
		return
	}
	if num, err = flags.GetInt("num"); err != nil {
		return
	}
	return
}
