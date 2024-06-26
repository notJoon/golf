package main

import (
	"go/token"
	"reflect"
	"testing"
)

func TestAnalyzeLifetime(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    *VariableTracker
		wantErr bool
	}{
		{
			name: "single declaration and use",
			src:  `package main; func main() { var x int; x = 5 }`,
			want: &VariableTracker{
				Decls: map[string]token.Pos{"x": 40},
				Uses:  map[string][]token.Pos{"x": {33, 40}},
			},
			wantErr: false,
		},
		{
			name:    "no declarations",
			src:     `package main; func main() {}`,
			want:    &VariableTracker{Decls: make(map[string]token.Pos), Uses: make(map[string][]token.Pos)},
			wantErr: false,
		},
		{
			name: "multiple declarations and uses",
			src:  `package main; func main() { var x int; var y int; x = 5; y = x }`,
			want: &VariableTracker{
				Decls: map[string]token.Pos{"x": 51, "y": 58},
				Uses:  map[string][]token.Pos{"x": {33, 51, 62}, "y": {44, 58}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AnalyzeLifetime(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeLifetime(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnalyzeLifetime(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
