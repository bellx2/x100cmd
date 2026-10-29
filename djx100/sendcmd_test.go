package djx100

import (
	"testing"
	"time"

	"go.bug.st/serial"
)

// 指定データを少しずつ返すテスト用ポート (未使用メソッドは埋め込みのまま)
type fakePort struct {
	serial.Port
	data  []byte
	chunk int
}

func (p *fakePort) Write(b []byte) (int, error) { return len(b), nil }

func (p *fakePort) Read(b []byte) (int, error) {
	n := min(p.chunk, len(p.data), len(b))
	copy(b, p.data[:n])
	p.data = p.data[n:]
	return n, nil
}

func TestSendCmd(t *testing.T) {
	// エコー行を読み飛ばして応答行を返す (分割受信も扱う)
	port := &fakePort{data: []byte("AL~VER\r\nver 4.00-001\r\n"), chunk: 3}
	got, err := SendCmd(port, "AL~VER")
	if err != nil || got != "ver 4.00-001" {
		t.Errorf("SendCmd() = %q, %v", got, err)
	}

	// 応答がなければタイムアウトする
	orig := cmdTimeout
	cmdTimeout = 50 * time.Millisecond
	defer func() { cmdTimeout = orig }()
	_, err = SendCmd(&fakePort{chunk: 1}, "AL~VER")
	if err == nil || err.Error() != "response timeout" {
		t.Errorf("SendCmd() error = %v, want response timeout", err)
	}
}
