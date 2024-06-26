package main

import (
	"fmt"
	"go/token"

	"github.com/notJoon/golf"
)

func main() {
	src := `
package main

func add(i, j int) int {
	return i + j
}

func main() { 
	var x int
	x = 5

	y := 10
	y += 1
	println(y)

	for i := 0; i < 5; i++ {
		println(i)
	}

	if x > 0 {
		println(add(x, y))
	}
}`
	fset := token.NewFileSet()
	tracker, err := golf.AnalyzeLifetime(src, fset)
	if err != nil {
		fmt.Println("Error analyzing source:", err)
		return
	}

	golf.PrintLifetime(tracker, fset)
}
