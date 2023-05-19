/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// bankCmd represents the bank command
var bankCmd = &cobra.Command{
	Use:   "bank",
	Short: "Bank Control",
	Args: cobra.MinimumNArgs(1),
}

var bankReadCmd = &cobra.Command{
	Use:   "read [A-Z]",
	Short: "Read Bank Name",
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dataOrg, err := djx100.ReadBankData(port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bankName := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		if len(args) > 0 {
			bankName = strings.ToUpper(args[0])
		}
		for _, v := range bankName {
			str, err := djx100.ParseBankName(dataOrg, string(v))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf(`"%s","%s"`+"\n", string(v), str)
		}
	},
}

var bankWriteCmd = &cobra.Command{
	Use:   "write [A-Z] [bank_name]",
	Short: "Write Bank Name",
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dataOrg, err := djx100.ReadBankData(port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bankName := strings.ToUpper(args[0])
		str, err := djx100.ParseBankName(dataOrg, bankName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Confirmation
		if cmd.Flag("yes").Value.String() != "true" {

			fmt.Printf("CH:%s %s -> %s\n", bankName, str, args[1])

			prompt := promptui.Prompt{
				Label:    "Write Bank Data",
				IsConfirm: true,
			}
			_, err = prompt.Run()
			if err != nil {
				os.Exit(1)	// No
			}	
		}else{
			fmt.Printf("Write : %s\n", args[1])
		}

		// Write

		n := args[1]
		if n == "NONE" {
			n = ""
		}
		dataNew, err := djx100.SetBankName(dataOrg, bankName, n)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := djx100.WriteBankData(port, dataNew)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(res)

		if cmd.Flag("restart").Value.String() == "true" {
			djx100.RestartCmd(port)
		}
	},
}

func init() {
	rootCmd.AddCommand(bankCmd)
	bankCmd.AddCommand(bankReadCmd)
	bankCmd.AddCommand(bankWriteCmd)
	bankWriteCmd.Flags().BoolP("yes", "y", false, "Without Confirmation")
	bankWriteCmd.Flags().BoolP("restart", "r", false, "Restart")
}
