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

// writeCmd represents the write command

var writeCmd = &cobra.Command{
	Use:   "write <channel>",
	Short: "Write Channel Data",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ch, _:= strconv.Atoi(args[0])
		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Read Current Channel Data
		data, err := djx100.ReadChData(port, ch)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if (cmd.Flag("debug").Value.String() == "true") {
			fmt.Println("Before_Data:",data)
		}
		chData, _ := djx100.ParseChData(data)
		chDataOrig := chData

		freq, _ := strconv.ParseFloat(cmd.Flag("freq").Value.String(), 64)
		if (freq != 0){
			chData.Freq = freq
		}

		mode := cmd.Flag("mode").Value.String()
		if (mode != ""){
			chData.Mode = djx100.ChMode2Num(mode)
			if (chData.Mode == -1){
				err := fmt.Errorf("invalid mode: %s", mode)
				fmt.Println(err)
				os.Exit(1)
			}
		}
		name := cmd.Flag("name").Value.String()
		if (name != ""){
			chData.Name = name
		}
		if (name == "NONE"){
			chData.Name = ""
		}

		if chData.IsEmpty() {
			err := fmt.Errorf("empty channel. freq Required")
			fmt.Println(err)
			os.Exit(1)
		}

		newData, err:= djx100.MakeChData(data, chData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if (cmd.Flag("debug").Value.String() == "true") {
			fmt.Println("After_Data:",newData)
		}

		// Confirmation
		if cmd.Flag("yes").Value.String() != "true" {

			fmt.Printf("Before: %s\n", chDataOrig.String())
			fmt.Printf("Afrer : %s\n", chData.String())

			prompt := promptui.Prompt{
				Label:    "Write Channel Data",
				IsConfirm: true,
			}
			_, err = prompt.Run()
			if err != nil {
				os.Exit(1)	// No
			}	
		}else{
			fmt.Printf("Write : %s\n", chData.String())
		}

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
	rootCmd.AddCommand(writeCmd)

	writeCmd.Flags().BoolP("restart", "r", false, "Send Restart Command")
	writeCmd.Flags().Float32P("freq", "f", 0, "Freqency")
	writeCmd.Flags().StringP("mode", "m", "", "Mode [FM,NFM,AM,NAM,T98,T102_B54...]")
	writeCmd.Flags().StringP("name", "n", "", "Name")
	writeCmd.Flags().BoolP("yes", "y", false, "Without Confirmation")
}
