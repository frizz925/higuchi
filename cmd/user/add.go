package user

import (
	"fmt"
	"syscall"

	"github.com/frizz925/higuchi/internal/auth"
	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var addCmd = &cobra.Command{
	Use:    "add",
	Short:  "Add user into passwords file",
	PreRun: argsCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCmd(cmd, args, runAdd)
	},
}

func runAdd(cmd *cobra.Command, args []string, h *hasher.MD5Hasher, users auth.Users) (auth.Users, error) {
	for _, user := range args {
		password, err := promptPassword(cmd, user)
		if err != nil {
			return nil, err
		}
		ad, err := hasher.NewMD5Digest(h, password)
		if err != nil {
			return nil, err
		}
		users[user] = ad
	}
	return users, nil
}

func promptPassword(cmd *cobra.Command, user string) (string, error) {
	out := cmd.OutOrStdout()
	for {
		fmt.Fprintf(out, "Enter password for %s: ", user)
		password, err := readPassword()
		if err != nil {
			return "", err
		}
		fmt.Fprintf(out, "\nConfirm password for %s: ", user)
		confirmation, err := readPassword()
		if err != nil {
			return "", err
		}
		if password == confirmation {
			fmt.Fprintf(out, "\nAdded new user %s\n", user)
			return password, nil
		}
		fmt.Fprintln(out, "\nPasswords mismatch. Please try again.")
	}
}

func readPassword() (string, error) {
	b, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
