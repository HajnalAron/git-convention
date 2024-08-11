package commands

import (
	"github.com/hajnalaron/git-convention-cli/internal/config"
	"github.com/spf13/cobra"
)

var configView = &cobra.Command{
	Use:   "config-view",
	Short: "Display the current configuration",
	Run:   displayConfig,
}

func displayConfig(cmd *cobra.Command, args []string) {
	config.ShowConfig(conf)
}
