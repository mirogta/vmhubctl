package cmd

import (
	"github.com/mirogta/vmhubctl/hub"
	"github.com/spf13/cobra"
)

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboots VM Hub 3",
	RunE:  reboot,
}

func init() {
	rootCmd.AddCommand(rebootCmd)
	// enableCmd.Flags().StringP("name", "n", "", "Set a name")
}

func reboot(cmd *cobra.Command, args []string) error {

	h := hub.NewHub()
	defer h.Logout()
	err := h.Login()
	if err != nil {
		return err
	}
	h.Reboot()
	return nil
}
