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
	d = append(d, ChData{Freq: 123.100, Mode: 2, Step: 3, Name: "0123456789012345678901234567", Att: 0, ShiftFreq: -3.1125, OffsetStep: true, Sq: 3, Tone: 4, DCS: 7})

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
}