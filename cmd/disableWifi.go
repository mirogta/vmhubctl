package cmd

import (
	"github.com/mirogta/vmhubctl/hub"
	"github.com/spf13/cobra"
)

var disableWifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "Disable Wifi Settings",
	RunE:  disableWifi,
}

func init() {
	disableCmd.AddCommand(disableWifiCmd)
}

func disableWifi(cmd *cobra.Command, args []string) error {
	h := hub.NewHub()
	defer h.Logout()
	err := h.Login()
	if err != nil {
		return err
	}
	h.SetWifi24GHzEnabled(false)
	h.SetWifi5GHzEnabled(false)
	return nil
}
