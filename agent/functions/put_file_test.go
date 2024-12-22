package functions

import (
	"os"
	"path/filepath"
	"testing"

	"github/clover0/github-issue-agent/store"
)

func TestPutFile(t *testing.T) {
	tests := []struct {
		name        string
		input       PutFileInput
		wantErr     bool
		checkFile   bool
		setupFunc   func(t *testing.T, dir string) // テスト前のセットアップ
		validateFunc func(t *testing.T, got store.File, err error)
	}{
		{
			name: "正常系：新規ファイル作成",
			input: PutFileInput{
				OutputPath:  "test.txt",
				ContentText: "test content",
			},
			checkFile: true,
			validateFunc: func(t *testing.T, got store.File, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.Path != "test.txt" || got.Content != "test content\n" {
					t.Errorf("unexpected result: got %+v", got)
				}
			},
		},
		{
			name: "正常系：ディレクトリ作成を伴うファイル作成",
			input: PutFileInput{
				OutputPath:  "new/dir/test.txt",
				ContentText: "nested content",
			},
			checkFile: true,
			validateFunc: func(t *testing.T, got store.File, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.Path != "new/dir/test.txt" || got.Content != "nested content\n" {
					t.Errorf("unexpected result: got %+v", got)
				}
			},
		},
		{
			name: "正常系：既存ファイルの上書き",
			input: PutFileInput{
				OutputPath:  "existing.txt",
				ContentText: "new content",
			},
			setupFunc: func(t *testing.T, dir string) {
				path := filepath.Join(dir, "existing.txt")
				if err := os.WriteFile(path, []byte("old content"), 0644); err != nil {
					t.Fatal(err)
				}
			},
			checkFile: true,
			validateFunc: func(t *testing.T, got store.File, err error) {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got.Path != "existing.txt" || got.Content != "new content\n" {
					t.Errorf("unexpected result: got %+v", got)
				}
			},
		},
		{
			name: "異常系：無効なパス（ファイル名なし）",
			input: PutFileInput{
				OutputPath:  "/invalid/path/",
				ContentText: "test content",
			},
			wantErr: true,
			validateFunc: func(t *testing.T, got store.File, err error) {
				if err == nil {
					t.Error("expected error but got nil")
				}
			},
		},
		{
			name: "異常系：書き込み権限なしディレクトリ",
			input: PutFileInput{
				OutputPath:  "noperm/test.txt",
				ContentText: "test content",
			},
			setupFunc: func(t *testing.T, dir string) {
				noPermDir := filepath.Join(dir, "noperm")
				if err := os.Mkdir(noPermDir, 0000); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: true,
			validateFunc: func(t *testing.T, got store.File, err error) {
				if err == nil {
					t.Error("expected error but got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用の一時ディレクトリを作成
			tmpDir := t.TempDir()
			
			// 現在のディレクトリを保存
			currentDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			
			// テストディレクトリに移動
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			
			// テスト終了時に元のディレクトリに戻る
			t.Cleanup(func() {
				os.Chdir(currentDir)
			})

			// セットアップ関数の実行
			if tt.setupFunc != nil {
				tt.setupFunc(t, tmpDir)
			}

			// テスト実行
			got, err := PutFile(tt.input)

			// バリデーション関数の実行
			if tt.validateFunc != nil {
				tt.validateFunc(t, got, err)
			}

			// ファイルの内容を確認
			if tt.checkFile && !tt.wantErr {
				content, err := os.ReadFile(tt.input.OutputPath)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
				}
				expectedContent := tt.input.ContentText
				if expectedContent[len(expectedContent)-1] != '\n' {
					expectedContent += "\n"
				}
				if string(content) != expectedContent {
					t.Errorf("file content mismatch: got %q, want %q", string(content), expectedContent)
				}
			}
		})
	}
}

func TestInitPutFileFunction(t *testing.T) {
	f := InitPutFileFunction()
	
	if f.Name != FuncPutFile {
		t.Errorf("unexpected function name: got %s, want %s", f.Name, FuncPutFile)
	}
	
	if f.Description != "Put new content to the file" {
		t.Errorf("unexpected description: got %s", f.Description)
	}
	
	if f.Func == nil {
		t.Error("function should not be nil")
	}
	
	// パラメータの検証
	params, ok := f.Parameters["properties"].(map[string]interface{})
	if !ok {
		t.Error("properties should be a map")
		return
	}
	
	// 必須パラメータの検証
	required, ok := f.Parameters["required"].([]string)
	if !ok {
		t.Error("required should be a string slice")
		return
	}
	
	expectedRequired := []string{"output_path", "content_text"}
	if len(required) != len(expectedRequired) {
		t.Errorf("unexpected required parameters: got %v, want %v", required, expectedRequired)
	}
}
