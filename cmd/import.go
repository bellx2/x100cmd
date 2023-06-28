/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <csv_filename>",
	Short: "import CSV data",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	
		r := csv.NewReader(file)
		record, err := r.Read()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		record[0] = strings.Trim(record[0], "\xef\xbb\xbf")
		cols := []string{"Channel","Freq","Mode","Step","Name"}
		for i, col := range cols {
			if (record[i] != col){
				fmt.Println("invalid csv format")
				os.Exit(1)
			}
		} 

		for {
			record, err := r.Read()
			if err != nil {
				fmt.Println(err)
				break
			}
			ch, _:= strconv.Atoi(record[0])
			data := djx100.BaseData
			if (cmd.Flag("overwrite").Value.String() == "false"){
				data, err = djx100.ReadChData(port, ch)
				if err != nil {
					fmt.Println(err)
					break
				}
			}

			chData, _ := djx100.ParseChData(data)
			freq, _ := strconv.ParseFloat(record[1], 64)
			if (freq == 0){
				// 0 の場合は消去実行
				res, err := djx100.WriteChData(port, ch, djx100.BaseData)
				if err != nil {
					fmt.Println(ch, err)
					continue
				}
				fmt.Printf("%d : %s : %s\n", ch, res, "clear")
				continue
			}else{
				chData.Freq = freq
			}
			mode := record[2]
			if (mode != ""){
				chData.Mode = djx100.ChMode2Num(mode)
				if (chData.Mode == -1){
					err := fmt.Errorf("invalid mode: %s", mode)
					fmt.Println(ch, err)
					continue
				}
			}
			step := record[3]
			if (step != ""){
				chData.Step = djx100.ChStep2Num(step)
				if (chData.Step == -1){
					err := fmt.Errorf("invalid step: %s", step)
					fmt.Println(ch, err)
					continue
				}
			}
			name := record[4]
			if (name != ""){
				chData.Name = name
			}
			if (name == "NONE"){
				chData.Name = ""
			}

			if len(record) > 5 {

				offset := record[5]
				if (offset != ""){
					if (offset == "OFF"){
						chData.OffsetStep = false
					}else if (offset == "ON"){
						chData.OffsetStep = true
					}
				}

				shift_freq, err := strconv.ParseFloat(record[6], 64)
				if err == nil {
					chData.ShiftFreq = shift_freq
				}

				att := record[7]
				if (att != ""){
					chData.Att = djx100.ChAtt2Num(att)
					if (chData.Att == -1){
						err := fmt.Errorf("invalid att: %s", att)
						fmt.Println(ch, err)
						continue
					}
				}

				sq := record[8]
				if (sq != ""){
					chData.Sq = djx100.ChSq2Num(sq)
					if (chData.Sq == -1){
						err := fmt.Errorf("invalid sq: %s", sq)
						fmt.Println(ch, err)
						continue
					}
				}

				tone := record[9]
				if (tone != ""){
					chData.Tone = djx100.ChTone2Num(tone)
					if (chData.Tone == -1){
						err := fmt.Errorf("invalid tone: %s", tone)
						fmt.Println(ch, err)
						continue
					}
				}

				dcs :=  fmt.Sprintf("%03s",record[10])
				if (dcs != "000"){
					chData.DCS = djx100.ChDCS2Num(dcs)
					if (chData.DCS == -1){
						err := fmt.Errorf("invalid dcs: %s", dcs)
						fmt.Println(ch, err)
						continue
					}
				}

				bank := record[11]
				if (bank != ""){
					chData.Bank = bank
					if bank == "NONE" {
						chData.Bank = ""
					}
				}

				// 拡張
				if len(record) > 12 {

					lat, err := strconv.ParseFloat(record[12], 64)
					if err == nil {
						chData.Lat = lat
					}

					lon, err := strconv.ParseFloat(record[13], 64)
					if err == nil {
						chData.Lon = lon
					}

					skip := record[14]
					if (skip != ""){
						if (skip == "OFF"){
							chData.Skip = false
						}else if (skip == "ON"){
							chData.Skip = true
						}
					}

					if len(record) > 15 {
						ext := record[15]
						if (ext != ""){
							chData.Ext = ext
						}
					}
				}
			}

			newData, err:= djx100.MakeChData(data, chData)
			if err != nil {
				fmt.Println(ch, err)
				continue
			}

			res, err := djx100.WriteChData(port, ch, newData)
			if err != nil {
				fmt.Println(ch, err)
				continue
			}
			if (cmd.Flag("verbose").Value.String() == "true") {
				fmt.Printf("%d : %s : %s\n", ch, res, chData.String())
			}else{
				fmt.Printf("%d[%s]: %.6f %s\n", ch, res, chData.Freq, chData.Name)
			}
		}
		
		if cmd.Flag("restart").Value.String() == "true" {
			err := djx100.RestartCmd(port)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	chCmd.AddCommand(importCmd)
	importCmd.Flags().BoolP("overwrite", "o", true, "Overwrite empty fields with default")
	importCmd.Flags().BoolP("restart", "r", false, "Send Restart Command")
	importCmd.Flags().BoolP("verbose", "v", false, "Make the operation more talkative")
}
