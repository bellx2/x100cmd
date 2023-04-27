/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear Channel Data",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch, _:= strconv.Atoi(args[0])
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flag("yes").Value.String() != "true" {
			fmt.Printf("Channel: %d\n", ch)
			prompt := promptui.Prompt{
				Label:    "Clear Channel Data",
				IsConfirm: true,
			}
			_, err = prompt.Run()
			if err != nil {
				os.Exit(1)	// No
			}	
		}
		newData := djx100.BaseData
		res, err := djx100.WriteChData(port, ch, newData)
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
	rootCmd.AddCommand(clearCmd)
	clearCmd.Flags().BoolP("yes", "y", false, "Without Confirmation")
	clearCmd.Flags().BoolP("restart", "r", false, "Send Restart Command")
}
