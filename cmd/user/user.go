package user

import (
	"fmt"
	"os"

	"github.com/frizz925/higuchi/internal/auth"
	"github.com/frizz925/higuchi/internal/config"
	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/spf13/cobra"
)

type userCommandFunc func(cmd *cobra.Command, args []string, h *hasher.Argon2Hasher, users auth.Argon2Users) (auth.Argon2Users, error)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage web proxy users",
}

func Command() *cobra.Command {
	return userCmd
}

func argsCheck(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		fmt.Fprintln(os.Stderr, "Provide at least one username in the arguments")
		os.Exit(1)
	}
}

func runCmd(cmd *cobra.Command, args []string, fn userCommandFunc) error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}
	pepper, err := cfg.Filters.Auth.Pepper()
	if err != nil {
		return fmt.Errorf("failed to decode pepper: %v", err)
	}
	h := hasher.NewArgon2Hasher(pepper)
	aa := auth.NewArgon2Auth(h)

	pf := cfg.Filters.Auth.PasswordsFile
	if _, err := os.Stat(pf); os.IsNotExist(err) {
		fmt.Println("Creating new passwords file...")
		if err := os.WriteFile(pf, nil, 0600); err != nil {
			return fmt.Errorf("failed to create passwords file: %v", err)
		}
	}
	users, err := aa.ReadPasswordsFile(pf)
	if err != nil {
		return fmt.Errorf("failed to read passwords file: %v", err)
	}

	newUsers, err := fn(cmd, args, h, users)
	if err != nil {
		return fmt.Errorf("failed to update users: %v", err)
	}

	err = aa.WritePasswordsFile(pf, newUsers)
	if err != nil {
		return fmt.Errorf("failed to write passwords file: %v", err)
	}
	return nil
}

func init() {
	userCmd.AddCommand(addCmd, delCmd)
}
