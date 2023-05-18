/*
Copyright © 2023 Ryu Tanabe <bellx2@gmali.com>
*/

package djx100

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type ChData struct {
	Freq float64
	Mode int
	Step int
	OffsetStep bool
	Name string
	ShiftFreq float64
	Att int
	Sq int
	Tone int
	DCS int
	Bank string
}
func (d ChData) IsEmpty() bool {
	return d.Freq == 0
}
func (d ChData) String() string {
	return fmt.Sprintf(`{"freq":%f, "mode":"%s", "step":"%s", "name":"%s", "offset":"%s", "shift_freq":"%f", "att":"%s", "sq":"%s", "tone":"%s", "dcs":"%s", "bank":"%s", "empty": %v}`, d.Freq, ChMode[d.Mode], ChStep[d.Step], d.Name, ChOffsetStep2Str(d.OffsetStep), d.ShiftFreq, ChAtt[d.Att], ChSq[d.Sq], ChTone[d.Tone], ChDCS[d.DCS],d.Bank,d.IsEmpty())
}

func(d *ChData) SetName(name string){
	n := strings.TrimRight(name, "\x00")
	d.Name = n
}

var ChMode = []string{"FM", "NFM", "AM", "NAM", "T98", "T102_B54", "DMR", "T61_typ1", "T61_typ2","T61_typ3","T61_typ4","", "", "dPMR","DSTAR","C4FM","AIS","ACARS","POCSAG","12KIF_W","12KIF_N" }

func ChMode2Num(mode string) (int){
	for i, v := range ChMode {
		if v == mode {
			return i
		}
	}
	return -1
}

var ChStep = []string{"1k","5k","6k25","8k33","10k","12k5","15k","20k","25k","30k","50k","100k","125k","200k"}

func ChStep2Num(step string) (int){
	for i, v := range ChStep {
		if v == step {
			return i
		}
	}
	return -1
}

func ChOffsetStep2Str(offset bool) (string){
	if offset {
		return "ON"
	}
	return "OFF"
}

var ChAtt = []string{"OFF","10db","20db"}
func ChAtt2Num(att string) (int){
	for i, v := range ChAtt {
		if v == att {
			return i
		}
	}
	return -1
}

var ChSq = []string{"OFF","CTCSS","DCS","R_CTCSS","R_DCS","JR","MSK"}
func ChSq2Num(sql string) (int){
	for i, v := range ChSq {
		if v == sql {
			return i
		}
	}
	return -1
}

var ChTone = []string{"670","693","719","744","770","797","825","854","885","915","948","974","1000","1035","1072","1109","1148","1188","1230","1273","1318","1355","1413","1462","1415","1567","1598","1622","1655","1679","1713","1738","1773","1799","1835","1862","1899","1928","1966","1995","2035","2065","2107","2181","2257","2291","2336","2418","2503","2541"}
func ChTone2Num(tone string) (int){
	for i, v := range ChTone {
		if v == tone {
			return i
		}
	}
	return -1
}

var ChDCS = []string{"017","023","025","026","031","032","036","043","047","050","051","053","054","065","071","072","073","074","114","115","116","122","125","131","132","134","143","145","152","155","156","162","165","172","174","205","212","223","225","226","243","244","245","246","251","252","255","261","263","265","266","271","274","306","311","315","325","331","332","343","346","351","356","364","365","371","411","412","413","423","431","432","445","446","452","454","455","462","464","465","466","503","506","516","523","526","532","546","565","606","612","624","627","631","632","654","662","664","703","712","723","731","732","734","743","754"}
func ChDCS2Num(dcs string) (int){
	for i, v := range ChDCS {
		if v == dcs {
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
	return "", errors.New("DJ-X100 not found")
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

	chByte, _ := hex.DecodeString(str)
	var freq uint32
	buf := bytes.NewBuffer(chByte[0x00:0x04])
	binary.Read(buf, binary.LittleEndian, &freq)
	d.Freq = float64(freq)/1000000
	d.Mode = int(chByte[0x04])
	d.Step = int(chByte[0x05])
	if chByte[0x06] == 0x01 {
		d.OffsetStep = true
	}else{
		d.OffsetStep = false
	}
	name := string(chByte[0x2b:0x48])
	name_utf8, _ := SJIStoUTF8(string(name))
	d.SetName(name_utf8) 

	var sfreq int32
	buf_s := bytes.NewBuffer(chByte[0x48:0x4c])
	binary.Read(buf_s, binary.LittleEndian, &sfreq)
	d.ShiftFreq = float64(sfreq)/1000000
	d.Att = int(chByte[0x4c])
	d.Sq = int(chByte[0x4d])
	d.Tone = int(chByte[0x4e])
	d.DCS = int(chByte[0x4f])

	bank_str := ""
	for i, v := range chByte[0x11:0x2B] {
		if(int(v) == 1){
			bank_str += fmt.Sprintf("%c", 0x41+i)
		}
	}
	d.Bank = bank_str

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

	buf := &bytes.Buffer{}
	_ = binary.Write(buf, binary.LittleEndian, int32(chData.Freq * 1000000))
	copy(chByte[0x00:0x04], buf.Bytes())	

	chByte[0x04] = byte(chData.Mode)
	chByte[0x05] = byte(chData.Step)
	if chData.OffsetStep {
		chByte[0x06] = 0x01
	}else{
		chByte[0x06] = 0x00
	}
	chByte[0x4c] = byte(chData.Att)
	chByte[0x4d] = byte(chData.Sq)
	chByte[0x4e] = byte(chData.Tone)
	chByte[0x4f] = byte(chData.DCS)

	buf_s := &bytes.Buffer{}
	_ = binary.Write(buf_s, binary.LittleEndian, int32(chData.ShiftFreq * 1000000))
	copy(chByte[0x48:0x4c], buf_s.Bytes())

	name_sjis, _ := UTF8toSJIS(chData.Name)
	for i := 0; i < 28; i++ {
		chByte[0x2b+i] = 0x00
		if i < len(name_sjis) {
			chByte[0x2b+i] = name_sjis[i]
		}
	}
	chByte[0x47] = 0x00

	for _, v := range chData.Bank {
		chByte[0x11+(v-0x41)] = byte(0x01)
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