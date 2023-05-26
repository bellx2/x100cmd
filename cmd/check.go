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

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check Serial Port And DX-J100 Connection",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("** scan ports **")
		err := djx100.ListPorts()
		if err != nil {
			panic(err)
		}
		fmt.Println("\n** check connection **")
		fmt.Printf("PortName: %s\n", rootCmd.PersistentFlags().Lookup("port").Value.String())
		portName, err := djx100.GetPortName(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("DJ-X100 PortName: %s\n", portName)

		port, err := djx100.Connect(portName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("\n** send device check command **")
		response, err := djx100.SendCmd(port, "AL~WHO")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(response)

		fmt.Println("\n** current version **")
		freq, err := djx100.SendCmd(port, "AL~VER")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(freq)

		if (cmd.Flag("restart").Value.String() == "true") {
			fmt.Println("\n** send restart command **")
			err := djx100.RestartCmd(port)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("restart", "r", false, "Send Restart Command")
}
