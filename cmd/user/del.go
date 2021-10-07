package user

import (
	"fmt"

	"github.com/frizz925/higuchi/internal/auth"
	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:    "del",
	Short:  "Delete user from passwords file",
	PreRun: argsCheck,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCmd(cmd, args, runDel)
	},
}

func runDel(cmd *cobra.Command, args []string, h *hasher.Argon2Hasher, users auth.Argon2Users) (auth.Argon2Users, error) {
	out := cmd.OutOrStdout()
	for _, user := range args {
		fmt.Fprintf(out, "Deleting user %s\n", user)
		delete(users, user)
	}
	return users, nil
}
