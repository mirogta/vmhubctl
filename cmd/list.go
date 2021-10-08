package cmd

import (
	"github.com/mirogta/vmhubctl/hub"
	"github.com/spf13/cobra"
)

// listCmd represents the say command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List a subset of VM Hub 3 settings",
	RunE:  List,
}

var listAll bool

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "list all discoverable settings")

	rootCmd.AddCommand(listCmd)
}

// List returns a subset of settings from the VM Hub
func List(cmd *cobra.Command, args []string) error {

	h := hub.NewHub()
	defer h.Logout()
	if err := h.Login(); err != nil {
		return err
	}

	if listAll {
		h.ListAll()
	} else {
		h.ListSelected()
	}
	return nil
}
