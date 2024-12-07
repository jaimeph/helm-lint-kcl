package validator

import (
	"os"
	"testing"
)

func LoadFileToBytes(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return data
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		values  []byte
		schemas []byte
		wantErr bool
	}{
		{
			name:    "valid input",
			values:  LoadFileToBytes("tests/test1.k"),
			schemas: LoadFileToBytes("tests/test1.yaml"),
			wantErr: false,
		},
		{
			name:    "invalid input",
			values:  []byte("invalid: value"),
			schemas: []byte("schema: definition"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.values, tt.schemas)
			if err := v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("'%s' - Validate() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
