/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/spf13/cobra"
)

// otpCmd represents the otp command
var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Auth OTP",
	Run: func(cmd *cobra.Command, args []string) {

		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := djx100.AuthOTP(port, rootCmd.PersistentFlags().Lookup("url").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(otpCmd)
	otpCmd.Flags().BoolP("verbose", "v", false, "Make the operation more talkative")
}
