/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// chCmd represents the ch command
var chCmd = &cobra.Command{
	Use:   "ch <command>",
	Short: "Channel Control",
}

func init() {
	rootCmd.AddCommand(chCmd)
}
