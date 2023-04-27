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

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Send Control Command",
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart",
	Run: func(cmd *cobra.Command, args []string) {
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = djx100.RestartCmd(port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("OK")
	},
}

var GPSCmd = &cobra.Command{
	Use:   "gps",
	Short: "GPS Info",
	Run: func(cmd *cobra.Command, args []string) {
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		response, err := djx100.SendCmd(port, "AL~GPS")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

var SQLCmd = &cobra.Command{
	Use:   "sql <level>",
	Short: "Set SQL [0-35]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		response, err := djx100.SendCmd(port, "AL~SQL" + args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

var VolCmd = &cobra.Command{
	Use:   "vol <level>",
	Short: "Set Volume [0-35]",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		response, err := djx100.SendCmd(port, "AL~VOL" + args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

var FreqCmd = &cobra.Command{
	Use:   "freq [MHz]",
	Short: "Get or Set current frequency",
	Run: func(cmd *cobra.Command, args []string) {
		freq := ""
		if len(args) != 0 {
			freq = args[0]
		}
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		response, err := djx100.SendCmd(port, "AL~FREQ"+freq)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.AddCommand(restartCmd)
	execCmd.AddCommand(GPSCmd)
	execCmd.AddCommand(SQLCmd)
	execCmd.AddCommand(VolCmd)
	execCmd.AddCommand(FreqCmd)
}
