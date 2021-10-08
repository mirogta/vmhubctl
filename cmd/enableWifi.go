package cmd

import (
	"github.com/mirogta/vmhubctl/hub"
	"github.com/spf13/cobra"
)

var enableWifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "Enable Wifi Settings",
	RunE:  enableWifi,
}

func init() {
	enableCmd.AddCommand(enableWifiCmd)
}

func enableWifi(cmd *cobra.Command, args []string) error {
	h := hub.NewHub()
	defer h.Logout()
	err := h.Login()
	if err != nil {
		return err
	}
	h.SetWifi24GHzEnabled(true)
	h.SetWifi5GHzEnabled(true)
	return nil
}
