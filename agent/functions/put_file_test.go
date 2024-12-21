package functions

import (
	"os"
	"testing"
)

func TestPutFile(t *testing.T) {
	tests := []struct {
		name        string
		input       PutFileInput
		wantErr     bool
		checkOutput func(t *testing.T, input PutFileInput)
	}{
		{
			name: "valid path and content",
			input: PutFileInput{
				OutputPath:  "testdata/testfile.txt",
				ContentText: "Hello, World!",
			},
			wantErr: false,
			checkOutput: func(t *testing.T, input PutFileInput) {
				content, err := os.ReadFile(input.OutputPath)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if string(content) != input.ContentText+"\n" {
					t.Errorf("got %s, want %s", content, input.ContentText+"\n")
				}
			},
		},
		{
			name: "invalid path",
			input: PutFileInput{
				OutputPath:  "/invalidpath/testfile.txt",
				ContentText: "Hello, World!",
			},
			wantErr: true,
			checkOutput: func(t *testing.T, input PutFileInput) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PutFile(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkOutput != nil {
				tt.checkOutput(t, tt.input)
			}
		})
	}
}
