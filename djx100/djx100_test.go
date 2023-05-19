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

	d = append(d, ChData{Freq: 433.100, Mode: 0, Step: 0, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: 0.0, OffsetStep: false, Sq: 0, Tone: 0, DCS: 0})
	d = append(d, ChData{Freq: 144.100, Mode: 1, Step: 1, Name: "0123456789012345678901234567", Att: 1, ShiftFreq: 1.0124, OffsetStep: true, Sq: 1, Tone: 3, DCS: 6})
	d = append(d, ChData{Freq: 123.100, Mode: 2, Step: 3, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: -3.1125, OffsetStep: true, Sq: 3, Tone: 4, DCS: 7, Bank: "ABCDXYZ"})
	d = append(d, ChData{Freq: 111.100, Mode: 4, Step: 3, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: 0.0, OffsetStep: true, Sq: 3, Tone: 4, DCS: 7, Bank: "ABCDE"})

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

		for k, v := range m {
			d := ChData{Freq: 433.100, Mode: 0, Step: 0, Name: k, Att: 0, ShiftFreq: 0.0, OffsetStep: false, Sq: 0, Tone: 0, DCS: 0}
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