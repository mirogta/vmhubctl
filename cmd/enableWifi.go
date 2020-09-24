package cmd

import (
	"github.com/mirogta/vmhubctl/hub"
	"github.com/spf13/cobra"
)

var enableWifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "Enable Wifi Settings",
	Run:   enableWifi,
}

func init() {
	enableCmd.AddCommand(enableWifiCmd)
}

func enableWifi(cmd *cobra.Command, args []string) {
	h := hub.NewHub()
	defer h.Logout()
	h.Login()
	h.SetWifi24GHzEnabled(true)
	h.SetWifi5GHzEnabled(true)
}
