/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <start> <end>",
	Short: "Convert mode values saved by V2 firmware to V4",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		start, err1 := strconv.Atoi(args[0])
		end, err2 := strconv.Atoi(args[1])
		if err1 != nil || err2 != nil || start < 0 || end > 999 || start > end {
			fmt.Println("invalid channel range (0-999)")
			os.Exit(1)
		}

		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		major, err := djx100.GetMajorVersion(port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if major < 4 {
			fmt.Println("convert is only for V4 firmware")
			os.Exit(1)
		}

		// 変換対象を先に集めて一覧表示する
		type target struct {
			ch   int
			data string
			mode int
		}
		targets := []target{}
		bar := pb.StartNew(end - start + 1)
		bar.SetMaxWidth(80)
		for ch := start; ch <= end; ch++ {
			data, err := djx100.ReadChData(port, ch)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			newData, ok := djx100.ConvertModeV2toV4(data)
			if ok {
				chData, _ := djx100.ParseChData(data)
				targets = append(targets, target{ch, newData, chData.Mode})
			}
			bar.Increment()
		}
		bar.Finish()

		if len(targets) == 0 {
			fmt.Println("no channels to convert")
			return
		}
		for _, t := range targets {
			fmt.Printf("%03d: %s (0x%02X -> 0x%02X)\n", t.ch, djx100.ChModeV2[t.mode], t.mode, t.mode+1)
		}

		// Confirmation
		if cmd.Flag("yes").Value.String() != "true" {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Convert %d channels", len(targets)),
				IsConfirm: true,
			}
			_, err = prompt.Run()
			if err != nil {
				os.Exit(1) // No
			}
		}

		for _, t := range targets {
			res, err := djx100.WriteChData(port, t.ch, t.data)
			if err != nil {
				fmt.Println(t.ch, err)
				continue
			}
			fmt.Printf("%03d : %s\n", t.ch, res)
		}
		if cmd.Flag("restart").Value.String() == "true" {
			djx100.RestartCmd(port)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	chCmd.AddCommand(convertCmd)
	convertCmd.Flags().BoolP("yes", "y", false, "Without Confirmation")
	convertCmd.Flags().BoolP("restart", "r", false, "Send Restart Command")
}
