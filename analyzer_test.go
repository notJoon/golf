package golf

import (
	"go/token"
	"testing"
)

func TestAnalyzeLifetime(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantVars int
		wantUses int
	}{
		{
			name: "Simple variable declaration and use",
			src: `
				package main
				func main() {
					x := 5
					println(x)
				}
			`,
			wantVars: 1,
			wantUses: 1,
		},
		{
			name: "Multiple variables and nested scope",
			src: `
				package main
				func main() {
					x := 5
					{
						y := x
						println(y)
					}
					println(x)
				}
			`,
			wantVars: 2,
			wantUses: 2,
		},
		{
			name: "Function with parameters",
			src: `
				package main
				func add(a, b int) int {
					return a + b
				}
				func main() {
					result := add(3, 4)
					println(result)
				}
			`,
			wantVars: 3,
			wantUses: 1,
		},
		{
			name: "For loop with range",
			src: `
				package main
				func main() {
					numbers := []int{1, 2, 3}
					for i, n := range numbers {
						println(i, n)
					}
				}
			`,
			wantVars: 3,
			wantUses: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			tracker, err := AnalyzeLifetime(tt.src, fset)
			if err != nil {
				t.Fatalf("AnalyzeLifetime(%s) error = %v", tt.name, err)
			}

			if got := len(tracker.Decls); got != tt.wantVars {
				t.Errorf("AnalyzeLifetime(%s) got %d variables, want %d", tt.name, got, tt.wantVars)
			}

			useCount := 0
			for _, uses := range tracker.Uses {
				useCount += len(uses)
			}
			if useCount != tt.wantUses {
				t.Errorf("AnalyzeLifetime(%s) got %d uses, want %d", tt.name, useCount, tt.wantUses)
			}

			if len(tracker.Scopes) == 0 {
				t.Errorf("AnalyzeLifetime(%s) did not create any scopes", tt.name)
			}
		})
	}
}
