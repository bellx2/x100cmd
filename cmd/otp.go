/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/bellx2/x100cmd/djx100"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

// otpCmd represents the otp command
var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "Auth OTP",
	Run: func(cmd *cobra.Command, args []string) {

		port, err := djx100.Connect(rootCmd.PersistentFlags().Lookup("port").Value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := authOTP(port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(res)
	},
}

// OTPサーバーURL (ビルド時に -ldflags "-X github.com/bellx2/x100cmd/cmd.otpURL=..." で埋め込む)
var otpURL string

// 外部サーバーへOTPを送信する前に同意を確認して認証する (同意しない場合は終了)
func authOTP(port serial.Port) (string, error) {
	server := rootCmd.PersistentFlags().Lookup("url").Value.String()
	if server == "" {
		server = otpURL
	}
	if server == "" {
		return "", errors.New("OTP server not configured (use --url)")
	}
	if rootCmd.PersistentFlags().Lookup("agree").Value.String() != "true" {
		host := server
		if u, err := url.Parse(server); err == nil && u.Host != "" {
			host = u.Host
		}
		fmt.Printf("OTP authentication sends the OTP read from the device to %s\n", host)
		// 対話できない場合は同意しなかったものとして扱う
		if fi, err := os.Stdin.Stat(); err != nil || fi.Mode()&os.ModeCharDevice == 0 {
			fmt.Println("not agreed (use --agree to send without confirmation)")
			os.Exit(1)
		}
		prompt := promptui.Prompt{
			Label:     "Send OTP to external server",
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err != nil {
			os.Exit(1) // No
		}
	}
	return djx100.AuthOTP(port, server)
}

func init() {
	rootCmd.AddCommand(otpCmd)
	otpCmd.Flags().BoolP("verbose", "v", false, "Make the operation more talkative")
}
