/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export <csv_filename>",
	Short: "export to CSV file",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Create(args[0])
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

		bomUtf8 := []byte{0xEF, 0xBB, 0xBF}	// UTF-8 BOM
		file.Write(bomUtf8)

		w := csv.NewWriter(file)
		w.Write([]string{"Channel","Freq","Mode","Step","Name","offset","shift_freq","att","sq","tone","dcs","bank"})

		bar := pb.StartNew(999)
		bar.SetMaxWidth(80)
		for ch:=1; ch<1000; ch++ {
			data, err := djx100.ReadChData(port, ch)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			chData, err := djx100.ParseChData(data)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if (cmd.Flag("all").Value.String() == "false" && chData.Freq == 0){
				bar.Increment()
				continue
			}
			w.Write([]string{fmt.Sprintf("%03d",ch), fmt.Sprintf("%.6f",chData.Freq), djx100.ChMode[chData.Mode], djx100.ChStep[chData.Step], chData.Name, djx100.ChOffsetStep2Str(chData.OffsetStep), fmt.Sprintf("%.6f",chData.ShiftFreq), djx100.ChAtt[chData.Att], djx100.ChSq[chData.Sq], djx100.ChTone[chData.Tone], djx100.ChDCS[chData.DCS], chData.Bank})
			bar.Increment()
		}
		bar.Finish()
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	chCmd.AddCommand(exportCmd)
	exportCmd.Flags().BoolP("all", "a", false, "Output All Channels")
}
