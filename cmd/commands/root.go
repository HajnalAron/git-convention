package commands

import (
	"fmt"
	"github.com/hajnalaron/git-convention-cli/internal/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	configPath string
	conf       *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "git-convention",
	Short: "A CLI tool for creating conventional branch names and commit messages",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to the configuration file")
	rootCmd.AddCommand(branch)
	rootCmd.AddCommand(commit)
	rootCmd.AddCommand(configView)
	cobra.OnInitialize(func() {
		confInUse, err := config.GetConfig(configPath)
		if err != nil {

		}
		conf = confInUse
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}
}
