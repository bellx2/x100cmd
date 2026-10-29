package djx100

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetToken(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		statusCode  int
		wantToken   string
		wantErr     bool
		errContains string
	}{
		{
			name:       "正常なトークン",
			response:   "1234567890abcdef",
			statusCode: http.StatusOK,
			wantToken:  "1234567890abcdef",
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
			response:    "NG",
			statusCode:  http.StatusBadRequest,
			wantErr:     true,
			errContains: "HTTP error: 400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			got, err := GetToken(server.URL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("GetToken() error = %v, want error containing %v", err, tt.errContains)
			}
			if got != tt.wantToken {
				t.Errorf("GetToken() = %v, want %v", got, tt.wantToken)
			}
		})
	}
}
