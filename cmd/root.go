/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "x100cmd",
	Short: "Unofficial Utility for ALINCO DJ-X100",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}


func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.Version="1.3.0"
	rootCmd.PersistentFlags().StringP("port", "p", "auto", "Serial Port Name")
	rootCmd.PersistentFlags().Bool("debug", false, "Show Debug Message")
}
