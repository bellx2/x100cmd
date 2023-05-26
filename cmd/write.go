/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// writeCmd represents the write command

var writeCmd = &cobra.Command{
	Use:   "write <channel_no>",
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
			fmt.Printf("Address: %05x\n", 0x20000 + (ch * 0x80))
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

		step := cmd.Flag("step").Value.String()
		if (step != ""){
			chData.Step = djx100.ChStep2Num(step)
			if (chData.Step == -1){
				err := fmt.Errorf("invalid step: %s", step)
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

		ostep := cmd.Flag("offset").Value.String()
		if (ostep != ""){
			if (ostep == "ON"){
				chData.OffsetStep = true
			}else{
				chData.OffsetStep = false
			}
		}

		if (cmd.Flag("shift_freq").Value.String() != ""){
			chData.ShiftFreq, _ = strconv.ParseFloat(cmd.Flag("shift_freq").Value.String(), 64)
		}

		att := cmd.Flag("att").Value.String()
		if (att != ""){
			chData.Att = djx100.ChAtt2Num(att)
			if (chData.Att == -1){
				err := fmt.Errorf("invalid att: %s", att)
				fmt.Println(err)
				os.Exit(1)
			}
		}

		sq := cmd.Flag("sq").Value.String()
		if (sq != ""){
			chData.Sq = djx100.ChSq2Num(sq)
			if (chData.Sq == -1){
				err := fmt.Errorf("invalid sq: %s", sq)
				fmt.Println(err)
				os.Exit(1)
			}
		}

		tone := cmd.Flag("tone").Value.String()
		if (tone != ""){
			chData.Tone = djx100.ChTone2Num(tone)
			if (chData.Tone == -1){
				err := fmt.Errorf("invalid tone: %s", tone)
				fmt.Println(err)
				os.Exit(1)
			}
		}
		
		dcs := cmd.Flag("dcs").Value.String()
		if (dcs != ""){
			chData.DCS = djx100.ChDCS2Num(dcs)
			if (chData.DCS == -1){
				err := fmt.Errorf("invalid dcs: %s", dcs)
				fmt.Println(err)
				os.Exit(1)
			}
		}

		bank := cmd.Flag("bank").Value.String()
		if (bank != ""){
			chData.Bank = strings.ToUpper(bank)
			if (bank == "NONE"){
				chData.Bank = ""
			}
		}

		skip := cmd.Flag("skip").Value.String()
		if (skip != ""){
			if (skip == "ON"){
				chData.Skip = true
			}else{
				chData.Skip = false
			}
		}

		if (cmd.Flag("lon").Value.String() != ""){
			chData.Lon, _ = strconv.ParseFloat(cmd.Flag("lon").Value.String(), 64)
		}

		if (cmd.Flag("lat").Value.String() != ""){
			chData.Lat, _ = strconv.ParseFloat(cmd.Flag("lat").Value.String(), 64)
		}

		ext := cmd.Flag("ext").Value.String()
		if (ext != ""){
			if(len(ext) != 96){
				err := fmt.Errorf("invalid ext (96chars): %s", ext)
				fmt.Println(err)
				os.Exit(1)
			}
			chData.Ext = ext
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
			fmt.Println("After_Data :",newData)
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
	chCmd.AddCommand(writeCmd)

	writeCmd.Flags().Float32P("freq", "f", 0, "Freqency")
	writeCmd.Flags().StringP("mode", "m", "", "Mode [FM,NFM,AM,NAM,T98,T102_B54...]")
	writeCmd.Flags().StringP("step", "s", "", "Step [6k25,10k,12k5,20k .... ]")
	writeCmd.Flags().StringP("name", "n", "", "Name")
	writeCmd.Flags().String("shift_freq", "", "Shift Freqency")
	writeCmd.Flags().String("offset", "", "OffestStep [ON,OFF]")
	writeCmd.Flags().String("att", "", "ATT [OFF,10db,20db]")
	writeCmd.Flags().String("sq", "", "Squelch [OFF,CTCSS,DCS,R_CTCSS,R_DCS,JR,MSK]")
	writeCmd.Flags().String("tone", "", "CTCSS Tone [670,693...2503,2541]")
	writeCmd.Flags().String("dcs", "", "DCS Code [017-754]")
	writeCmd.Flags().String("bank", "", "Bank [A-Z] ex. ABCDEZ")
	writeCmd.Flags().String("skip", "", "Skip [ON,OFF]")
	writeCmd.Flags().String("lat", "", "Latitude")
	writeCmd.Flags().String("lon", "", "Longitude")
	writeCmd.Flags().String("ext", "", "ExtData(0x50-0x7F) 96chars")
	writeCmd.Flags().BoolP("yes", "y", false, "Without Confirmation")
	writeCmd.Flags().BoolP("restart", "r", false, "Restart")
}
