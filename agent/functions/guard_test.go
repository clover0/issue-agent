package functions

import (
	"strings"
	"testing"
)

func TestGuardPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "正常系: 通常のパス",
			path:    "path/to/file.txt",
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "正常系: 空のパス",
			path:    "",
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "異常系: 親ディレクトリ参照",
			path:    "../path/to/file",
			wantErr: true,
			errMsg:  "contains not allowed '..'",
		},
		{
			name:    "異常系: チルダ使用",
			path:    "~/path/to/file",
			wantErr: true,
			errMsg:  "contains not allowed '~'",
		},
		{
			name:    "異常系: 重複スラッシュ",
			path:    "path//to/file",
			wantErr: true,
			errMsg:  "contains not allowed '//'",
		},
		{
			name:    "異常系: ルートパス開始",
			path:    "/path/to/file",
			wantErr: true,
			errMsg:  "starts with '/', not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := guardPath(tt.path)
			
			// エラーの有無を検証
			if (err != nil) != tt.wantErr {
				t.Errorf("guardPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーメッセージの検証
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("guardPath() error message = %v, want containing %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}
