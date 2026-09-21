package cmd

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		statusCode  int
		delay       time.Duration
		wantToken   string
		wantErr     bool
		errContains string
	}{
		{
			name:        "正常なトークン",
			response:    "1234567890abcdef",
			statusCode:  http.StatusOK,
			wantToken:   "1234567890abcdef",
			wantErr:     false,
		},
		{
			name:        "無効な長さのトークン",
			response:    "123",
			statusCode:  http.StatusOK,
			wantErr:     true,
			errContains: "token too short",
		},
		{
			name:        "無効な16進数",
			response:    "invalid hex token",
			statusCode:  http.StatusOK,
			wantErr:     true,
			errContains: "invalid hex token",
		},
		{
			name:        "HTTPエラー",
			response:    "error",
			statusCode:  http.StatusInternalServerError,
			wantErr:     true,
			errContains: "HTTP error: 500",
		},
		{
			name:        "タイムアウト",
			response:    "1234567890abcdef",
			delay:       2 * time.Second,
			statusCode:  http.StatusOK,
			wantErr:     true,
			errContains: "timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.delay > 0 {
					time.Sleep(tt.delay)
				}
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// テスト用にタイムアウトを短く設定
			defaultTimeout := 1 * time.Second
			http.DefaultClient.Timeout = defaultTimeout
			http.DefaultClient.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			got, err := GetToken(server.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("GetToken() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if got != tt.wantToken {
				t.Errorf("GetToken() = %v, want %v", got, tt.wantToken)
			}
		})
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}