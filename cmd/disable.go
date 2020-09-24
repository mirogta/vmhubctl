package cmd

import (
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable VM Hub 3 settings",
}

func init() {
	rootCmd.AddCommand(disableCmd)
}
