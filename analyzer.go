package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// VariableTracker keeps track of the declare and use of variables
type VariableTracker struct {
	Decls map[string]token.Pos
	Uses  map[string][]token.Pos
}

func NewVariableTracker() *VariableTracker {
	return &VariableTracker{
		Decls: make(map[string]token.Pos),
		Uses:  make(map[string][]token.Pos),
	}
}

func AnalyzeLifetime(src string) (*VariableTracker, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}

	tracker := NewVariableTracker()

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					tracker.Decls[ident.Name] = ident.Pos()
				}
			}
		case *ast.Ident:
			if x.Obj != nil && x.Obj.Kind == ast.Var {
				tracker.Uses[x.Name] = append(tracker.Uses[x.Name], x.Pos())
			}
		}
		return true
	})

	return tracker, nil
}

func PrintLifetime(tracker *VariableTracker, fset *token.FileSet) {
	for varName, declPos := range tracker.Decls {
		fmt.Printf("Variable declared: %s at %v\n", varName, fset.Position(declPos))
		if uses, ok := tracker.Uses[varName]; ok {
			for _, usePos := range uses {
				fmt.Printf("  Used at: %v\n", fset.Position(usePos))
			}
		}
	}
}
