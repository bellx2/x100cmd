/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read <channel>",
	Short: "Read Channel Data",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch, _:= strconv.Atoi(args[0])
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data, err := djx100.ReadChData(port, ch)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if (cmd.Flag("debug").Value.String() == "true") {
			fmt.Println("data:",data)
		}
		chData, err := djx100.ParseChData(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(chData.String())
		djx100.Close(port)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolP("debug", "d", false, "Show Debug Mode")
}
