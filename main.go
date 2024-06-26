package main

import (
	"go/parser"
	"go/token"
)

func main() {
	src := `
    package main

    func example() {
        x := 10
        {
            y := x
            fmt.Println(y)
        }
        fmt.Println(x)
    }
    `

	tracker, err := AnalyzeLifetime(src)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	parser.ParseFile(fset, "", src, 0)

	PrintLifetime(tracker, fset)
}
