package cmd

import (
	"github.com/mirogta/vmhubctl/hub"
	"github.com/spf13/cobra"
)

var disableWifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "Disable Wifi Settings",
	Run:   disableWifi,
}

func init() {
	disableCmd.AddCommand(disableWifiCmd)
}

func disableWifi(cmd *cobra.Command, args []string) {
	h := hub.NewHub()
	defer h.Logout()
	h.Login()
	h.SetWifi24GHzEnabled(false)
	h.SetWifi5GHzEnabled(false)
}
