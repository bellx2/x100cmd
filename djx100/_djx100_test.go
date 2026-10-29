package djx100

import (
    "testing"
)

// ... (前回までのテストケースは同じ) ...

// Bool2Strのテスト
func TestBool2Str(t *testing.T) {
    tests := []struct {
        input    bool
        expected string
        desc     string
    }{
        {true, "ON", "真値をONに変換"},
        {false, "OFF", "偽値をOFFに変換"},
    }

    for _, tt := range tests {
        t.Run(tt.desc, func(t *testing.T) {
            result := Bool2Str(tt.input)
            if result != tt.expected {
                t.Errorf("Bool2Str(%v) = %s; want %s (%s)", tt.input, result, tt.expected, tt.desc)
            } else {
                t.Logf("Bool2Str Test Success: %v -> %s (%s)", tt.input, result, tt.desc)
            }
        })
    }
}

// SJIS変換のテスト
func TestSJISMultiCheck(t *testing.T) {
    tests := []struct {
        input    byte
        expected bool
        desc     string
    }{
        {0x81, true, "第1水準漢字開始"},
        {0x9F, true, "第1水準漢字終了"},
        {0xE0, true, "第2水準漢字開始"},
        {0xFC, true, "第2水準漢字終了"},
        {0x80, false, "制御コード"},
        {0xA0, false, "1バイト文字"},
        {0x40, false, "ASCII文字"},
    }

    for _, tt := range tests {
        t.Run(tt.desc, func(t *testing.T) {
            result := SJISMultiCheck(tt.input)
            if result != tt.expected {
                t.Errorf("SJISMultiCheck(0x%02X) = %v; want %v (%s)", tt.input, result, tt.expected, tt.desc)
            } else {
                t.Logf("SJISMultiCheck Test Success: 0x%02X -> %v (%s)", tt.input, result, tt.desc)
            }
        })
    }
}

// チャンネルデータの空チェック
func TestIsEmpty(t *testing.T) {
    tests := []struct {
        data     ChData
        expected bool
        desc     string
    }{
        {ChData{}, true, "完全な空構造体"},
        {ChData{Freq: 0}, true, "周波数0のみ"},
        {ChData{Freq: 433.100}, false, "有効な周波数"},
        {ChData{Freq: 0, Name: "Test"}, true, "周波数0で名前あり"},
        {ChData{Freq: 0, Mode: 1}, true, "周波数0でモードあり"},
    }

    for _, tt := range tests {
        t.Run(tt.desc, func(t *testing.T) {
            result := tt.data.IsEmpty()
            if result != tt.expected {
                t.Errorf("IsEmpty() for %s = %v; want %v", tt.desc, result, tt.expected)
            } else {
                t.Logf("IsEmpty Test Success: %s -> %v", tt.desc, result)
            }
        })
    }
}

