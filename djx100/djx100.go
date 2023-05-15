/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package djx100

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var FreqMin = 20.0
var FreqMax = 470.0

type ChData struct {
	Freq float64
	Mode int
	Name string
	Step int
}
func (d ChData) IsEmpty() bool {
	return d.Freq == 0
}
func (d ChData) String() string {
	return fmt.Sprintf(`{"freq":%f, "mode":"%s", "step":"%s", "name":"%s", "empty": %v}`, d.Freq, ChMode[d.Mode], ChStep[d.Step], d.Name, d.IsEmpty())
}

var ChMode = []string{"FM", "NFM", "AM", "NAM", "T98", "T102_B54", "DMR", "T61_typ1", "T61_typ2","T61_typ3","T61_typ4","", "", "dPMR","DSTAR","C4FM","AIS","ACARS","POCSAG","12KIF_W","12KIF_N" }

func ChMode2Num(mode string) (int){
	if mode == "" {
		return -1
	}
	for i, v := range ChMode {
		if v == mode {
			return i
		}
	}
	return -1
}

var ChStep = []string{"1k","5k","6k25","8k33","10k","12k5","15k","20k","25k","30k","50k","100k","125k","200k"}

func ChStep2Num(step string) (int){
	if step == "" {
		return -1
	}
	for i, v := range ChStep {
		if v == step {
			return i
		}
	}
	return -1
}

var BaseData = "FFFFFFFF000700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000E4000000E400000000000000000000000180018001800180010000800100008001000080000080008000807B1700"

// シリアルポート一覧取得
func ListPorts() (error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return err
	}
	for _, port := range ports {
		if port.IsUSB {
			dev := "Unknown"
			if (port.VID == "3614" && port.PID == "D001"){
				dev = "DJ-X100!"
			}
			fmt.Printf("%s [%s:%s] %s\n",port.Name, port.VID, port.PID, dev)
		}else{
			fmt.Printf("%s\n",port.Name)
		}
	}
	return nil
}

// DJ-X100ポート名取得
func GetPortName(portName string) (string, error) {
	if portName != "auto" {
		return portName, nil
	}
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", err
	}
	for _, port := range ports {
		if port.IsUSB && port.VID == "3614" && port.PID == "D001" {
			return port.Name, nil
		}
	}
	return "", errors.New("DJ-X100 is not found")
}

// シリアルポート接続
func Connect(portName string) (serial.Port, error){
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	p, err := GetPortName(portName)
	if err != nil {
		return nil, err
	}
	port, err := serial.Open(p, mode)
	if err != nil {
		return nil, err
	}
	response, err := SendCmd(port, "AL~DJ-X100")
	if (err != nil || response != "OK"){
		return nil, errors.New("DJ-X100 is not conntected")
	}
	return port, nil
}

// シリアルポート切断
func Close(port serial.Port) (error){
	port.SetDTR(false)
	port.SetRTS(false)
	port.Close()
	return nil
}

// チャンネルデータ読み込み
func ReadChData(port serial.Port, ch int) (string, error){
	if ch < 0 || ch > 999 {
		return "", errors.New("channel number is out of range")
	}
	str, err := SendCmd(port, fmt.Sprintf("AL~F%05xM",0x20000 + (ch * 0x80)))
	if err != nil {
		return "", errors.New("read Command Error")
	}
	return str[0:0x100], nil
}

// チャンネルデータ解析
func ParseChData(str string)(ChData, error){

	if str[0:8] == "FFFFFFFF" {
		return ChData{}, nil
	}

	d := ChData{}

	freq, _ := strconv.ParseUint(str[6:8] + str[4:6] + str[2:4] + str[0:2], 16, 32)
	fmode, _ := strconv.ParseUint(str[8:10], 16, 8)
	fstep, _ := strconv.ParseUint(str[10:12], 16, 8)
	name, _ := hex.DecodeString(string(str[86:145]))	// SJIS
	name_utf8, _ := SJIStoUTF8(string(name))

	d.Freq = float64(freq)/1000000
	d.Mode = int(fmode)
	d.Step = int(fstep)
	d.Name = strings.TrimRight(name_utf8, "\x00")
	return d, nil
}

// チャンネルデータ作成
func MakeChData(dataOrg string, chData ChData) (string, error){

	if (dataOrg[0:8] == "FFFFFFFF" && chData.Freq == 0){
		return "", errors.New("empty channel. freq Required")
	}
	if dataOrg[0:8] == "FFFFFFFF" {
		dataOrg = BaseData
	}
	chByte, _ := hex.DecodeString(dataOrg)

	if (FreqMin < chData.Freq && chData.Freq < FreqMax){
		w_freq := int(chData.Freq * 1000000)
		w_freq_hex, _ := hex.DecodeString(fmt.Sprintf("%08X", w_freq))
		// fmt.Printf("freq: %d\n", w_freq)
		// fmt.Printf("freq_str: %x\n", w_freq_hex)
		chByte[0] = w_freq_hex[3]
		chByte[1] = w_freq_hex[2]
		chByte[2] = w_freq_hex[1]
		chByte[3] = w_freq_hex[0]			
	}	

	chByte[4] = byte(chData.Mode)
	chByte[5] = byte(chData.Step)

	name_sjis, _ := UTF8toSJIS(chData.Name)
	for i := 0; i < len(name_sjis); i++ {
		if i >= 30 {
			break
		}
		chByte[43 + i] = name_sjis[i]
	}
	for i := len(name_sjis); i<30 ; i++ {
		chByte[43 + i] = 0x00
	}

	return hex.EncodeToString(chByte), nil
}

// チャンネルデータ読み込み
func WriteChData(port serial.Port, ch int, data string) (string, error){
	if ch < 0 || ch > 999 {
		return "", errors.New("channel number is out of range")
	}
	str, err := SendCmd(port, fmt.Sprintf("AL~F%05xW%s",0x20000 + (ch * 0x80),data))
	if err != nil {
		return "", errors.New("write Command Error")
	}
	return str, nil
}

// コマンド送信
func SendCmd(port serial.Port , cmd string) (string, error){
	_, err := port.Write([]byte(cmd + "\r\n"))
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(port)
	scanner.Scan()	// skip first line
	scanner.Scan()	// response
	response := scanner.Text()
	return response, nil
}

// リスタート送信
func RestartCmd(port serial.Port) (error){
	_, err := port.Write([]byte("AL~RESTART\r\n"))
	if err != nil {
		return err
	}
	port.SetDTR(false)
	port.SetRTS(false)
	port.Close()
	return nil
}

// UTF-8 から ShiftJIS
func UTF8toSJIS(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
					return "", err
	}
	return string(ret), err
}

// ShiftJIS から UTF-8
func SJIStoUTF8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
					return "", err
	}
	return string(ret), err
}