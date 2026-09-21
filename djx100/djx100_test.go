package djx100

import (
	"testing"
)

func SetAndParse(CjData ChData) (ChData, error){
	d, err := MakeChData(BaseData, CjData)
	if err != nil {
		return ChData{}, err
	}
	return ParseChData(d)
}

func TestSetAndParse(t *testing.T){
	d := []ChData{}
	
	// テストデータ

	d = append(d, ChData{Freq: 433.100, Mode: 0, Step: 0, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: 0.0, OffsetStep: false, Sq: 0, Tone: 0, DCS: 0,
	Ext: "0100e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1700"})

	d = append(d, ChData{Freq: 144.100, Mode: 1, Step: 1, Name: "0123456789012345678901234567", Att: 1, ShiftFreq: 1.0124, OffsetStep: true, Sq: 1, Tone: 3, DCS: 6, 
	Ext:"0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1701"})
	
	d = append(d, ChData{Freq: 123.100, Mode: 2, Step: 3, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: -3.1125, 			OffsetStep: true, Sq: 3, Tone: 4, DCS: 7, Bank: "ABCDXYZ", Lat:35.681382, Lon:139.766084, Skip: true,
	Ext: "0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1703"})
	
	d = append(d, ChData{Freq: 111.100, Mode: 4, Step: 3, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: 1.12345, OffsetStep: true, Sq: 3, Tone: 4, DCS: 7, Bank: "ABCDE", Lat:35.681382, Lon:139.766084,
	Ext: "0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080408000807b1700"})

	d = append(d, ChData{Freq: 128.325, Mode: 0, Step: 3, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: 0.12222, OffsetStep: true, Sq: 3, Tone: 4, DCS: 7, Bank: "ABCDE", Lat:35.681382, Lon:139.766084,
	Ext: "0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080408000807b1700"})

	for _, v := range d {
		t.Run("SetAndParse", func(t *testing.T) {
			got, err := SetAndParse(v)
			if err != nil {
				t.Errorf("SetAndParse Error: %v", err)
			}
			if got != v {
				t.Errorf("SetAndParse MisMatch\nSet:%v\nRes:%v", v, got)
			}else{
				t.Logf("SetAndParse: %v", v)
			}
		})
	}

	t.Run("Name Max Check", func(t *testing.T) {

		m := map[string]string{}
		m["0123456789012345678901234567890123"] = "0123456789012345678901234567"
		m["アイウエオかきくけこさしすせそ"] = "アイウエオかきくけこさしすせ"
		m["0アイウエオかきくけこさしすせそ"] = "0アイウエオかきくけこさしす"
		m["0アイウエオかきくけこさしす0せそ"] = "0アイウエオかきくけこさしす0"
		m["アイウ"] = "アイウ"
		m["0アイウ"] = "0アイウ"
		m["G"] = "G"
		m["G1"] = "G1"
		m["G12"] = "G12"

		for k, v := range m {
			d := ChData{Freq: 433.100, Mode: 0, Step: 0, Name: k, Att: 0, ShiftFreq: 0.0, OffsetStep: false, Sq: 0, Tone: 0, DCS: 0}
			d.Ext = "0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1700"
			d_res := d
			d_res.Name = v
			got, err := SetAndParse(d)
			if err != nil {
				t.Errorf("SetAndParse Error: %v", err)
			}
			if got.String() != d_res.String() {	//他のパラメータ含めて比較
				if got.Name != d_res.Name {
					t.Errorf("Name Error %s -> %s \n", d_res.Name, got.Name)
				}else{
					t.Errorf("Data Error %s -> %s \n", d_res.String(), got.String())
				}
			}else{
				t.Logf("SetAndParse: %s -> %s", k, got.Name)
			}
		}
	})

	t.Run("Bank Check", func(t *testing.T) {

		m := map[string]string{}
		m["ABCDEFXYZ"] = "ABCDEFXYZ"
		m["ABCDEFXYZZZ"] = "ABCDEFXYZ"
		m["AAAAAABCDEFFFFXYZZZ"] = "ABCDEFXYZ"
		m["abc"] = "ABC"
		m["abc@%()"] = "ABC"

		for k, v := range m {
			d := ChData{Freq: 433.100, Mode: 0, Step: 0, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: 0.0, OffsetStep: false, Sq: 0, Tone: 0, DCS: 0, Bank: k}
			d.Ext = "0000e4000000e400000000000000000000000180018001800180010000800100008001000080000080008000807b1700"
			d_res := d
			d_res.Bank = v
			got, err := SetAndParse(d)
			if err != nil {
				t.Errorf("SetAndParse Error: %v", err)
			}
			if got.String() != d_res.String() {	//他のパラメータ含めて比較
				if got.Bank != d_res.Bank {
					t.Errorf("Bank Error %s -> %s \n", d_res.Bank, got.Bank)
				}else{
					t.Errorf("Data Error %s -> %s \n", d_res.String(), got.String())
				}
			}else{
				t.Logf("SetAndParse: %s -> %s", k, got.Bank)
			}
		}
	})

}

// ChMode2Numのテスト
func TestChMode2Num(t *testing.T) {
    tests := []struct {
        input    string
        expected int
        desc     string
    }{
        {"FM", 0, "基本モード：FM"},
        {"AM", 2, "基本モード：AM"},
        {"DMR", 6, "デジタルモード：DMR"},
        {"12KIF_N", 20, "最後のモード"},
        {"INVALID", -1, "無効なモード"},
        {"", -1, "空文字列"},
    }

    for _, tt := range tests {
        t.Run(tt.desc, func(t *testing.T) {
            result := ChMode2Num(tt.input)
            if result != tt.expected {
                t.Errorf("ChMode2Num(%q) = %d; want %d (%s)", tt.input, result, tt.expected, tt.desc)
            } else {
                t.Logf("ChMode2Num Test Success: %q -> %d (%s)", tt.input, result, tt.desc)
            }
        })
    }
}

// ChMode2Strのテスト
func TestChMode2Str(t *testing.T) {
    tests := []struct {
        input    int
        expected string
        desc     string
    }{
        {0, "FM", "基本モード：FM"},
        {2, "AM", "基本モード：AM"},
        {6, "DMR", "デジタルモード：DMR"},
        {20, "12KIF_N", "最後のモード"},
        {-1, "Unknown", "範囲外：負数"},
        {21, "Unknown", "範囲外：上限超過"},
    }

    for _, tt := range tests {
        t.Run(tt.desc, func(t *testing.T) {
            result := ChMode2Str(tt.input)
            if result != tt.expected {
                t.Errorf("ChMode2Str(%d) = %q; want %q (%s)", tt.input, result, tt.expected, tt.desc)
            } else {
                t.Logf("ChMode2Str Test Success: %d -> %q (%s)", tt.input, result, tt.desc)
            }
        })
    }
}
