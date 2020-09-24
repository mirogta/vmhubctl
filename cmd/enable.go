package cmd

import (
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable VM Hub 3 settings",
}

func init() {
	rootCmd.AddCommand(enableCmd)
}
