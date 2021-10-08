package cmd

import (
	"github.com/frizz925/higuchi/cmd/user"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "higuchi",
	Short: "Higuchi web proxy",
	Long:  `Higuchi is a performant and module web proxy written in Golang.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd, user.Command())
}
