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
			data, err := djx100.ReadChData(port, ch)
			if err != nil {
				fmt.Println(err)
				break
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
			fmt.Printf("%d : %s : %s\n", ch, res, chData.String())
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
	importCmd.Flags().BoolP("restart", "r", false, "Send Restart Command")
}
