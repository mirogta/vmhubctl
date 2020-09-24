package cmd

import (
	"github.com/mirogta/vmhubctl/hub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the say command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List a subset of VM Hub 3 settings",
	Run:   List,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// List returns a subset of settings from the VM Hub
func List(cmd *cobra.Command, args []string) {

	h := hub.NewHub()
	defer h.Logout()
	if err := h.Login(); err != nil {
		log.Error(err)
		return
	}

	h.ListAll()
}
