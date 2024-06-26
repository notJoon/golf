package golf

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Scope struct {
	ID     int
	Parent *Scope
	Vars   map[string]token.Pos
}

// VariableTracker keeps track of the declare and use of variables
type VariableTracker struct {
	Decls        map[string]token.Pos   // Variable declaration positions
	Uses         map[string][]token.Pos // Variable use positions
	Scopes       map[int]*Scope
	CurrentScope *Scope
	NextScopeID  int
}

func NewVariableTracker() *VariableTracker {
	globalScope := &Scope{ID: 0, Parent: nil, Vars: make(map[string]token.Pos)}
	return &VariableTracker{
		Decls:        make(map[string]token.Pos),
		Uses:         make(map[string][]token.Pos),
		Scopes:       map[int]*Scope{0: globalScope},
		CurrentScope: globalScope,
		NextScopeID:  1,
	}
}

func (vt *VariableTracker) PushScope() {
	newScope := &Scope{
		ID:     vt.NextScopeID,
		Parent: vt.CurrentScope,
		Vars:   make(map[string]token.Pos),
	}

	vt.Scopes[newScope.ID] = newScope
	vt.CurrentScope = newScope
	vt.NextScopeID++
}

func (vt *VariableTracker) PopScope() {
	if vt.CurrentScope.Parent != nil {
		vt.CurrentScope = vt.CurrentScope.Parent
	}
}

func (vt *VariableTracker) DeclareVar(name string, pos token.Pos) {
	vt.Decls[name] = pos
	vt.CurrentScope.Vars[name] = pos
}

func (vt *VariableTracker) UseVar(name string, pos token.Pos) {
	vt.Uses[name] = append(vt.Uses[name], pos)
}

func AnalyzeLifetime(src string, fset *token.FileSet) (*VariableTracker, error) {
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, err
	}

	tracker := NewVariableTracker()

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			tracker.PushScope()
			// Add function parameters as variables
			if x.Type.Params != nil {
				for _, field := range x.Type.Params.List {
					for _, name := range field.Names {
						tracker.DeclareVar(name.Name, name.Pos())
					}
				}
			}
			defer tracker.PopScope()
		case *ast.BlockStmt:
			tracker.PushScope()
			defer tracker.PopScope()
		case *ast.RangeStmt:
			tracker.PushScope()
			// Add range variables
			if x.Key != nil {
				if ident, ok := x.Key.(*ast.Ident); ok {
					tracker.DeclareVar(ident.Name, ident.Pos())
				}
			}
			if x.Value != nil {
				if ident, ok := x.Value.(*ast.Ident); ok {
					tracker.DeclareVar(ident.Name, ident.Pos())
				}
			}
			defer tracker.PopScope()
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					if x.Tok == token.DEFINE {
						tracker.DeclareVar(ident.Name, ident.Pos())
					} else {
						tracker.UseVar(ident.Name, ident.Pos())
					}
				}
			}
		case *ast.CallExpr:
			if ident, ok := x.Fun.(*ast.Ident); ok {
				// TOOD: check for builtin functions using a map
				if ident.Name == "println" || ident.Name == "print" {
					for _, arg := range x.Args {
						ast.Inspect(arg, func(n ast.Node) bool {
							if ident, ok := n.(*ast.Ident); ok {
								if ident.Obj != nil && ident.Obj.Kind == ast.Var {
									tracker.UseVar(ident.Name, ident.Pos())
								}
							}
							return true
						})
					}
				}
			}
		case *ast.Ident:
			// Only count as use if it's not part of a declaration
			if x.Obj != nil && x.Obj.Kind == ast.Var && !isPartOfDeclaration(x) {
				tracker.UseVar(x.Name, x.Pos())
			}
		}
		return true
	})

	return tracker, nil
}

func isPartOfDeclaration(ident *ast.Ident) bool {
	if ident.Obj == nil {
		return false
	}

	switch ident.Obj.Decl.(type) {
	case *ast.AssignStmt, *ast.ValueSpec, *ast.Field:
		return true
	}

	return false
}

func PrintLifetime(tracker *VariableTracker, fset *token.FileSet) {
	for scopeID, scope := range tracker.Scopes {
		fmt.Printf("Scope %d:\n", scopeID)
		if scope.Parent != nil {
			fmt.Printf("  Parent Scope: %d\n", scope.Parent.ID)
		}
		for varName, pos := range scope.Vars {
			fmt.Printf("  Variable declared: %s at %v\n", varName, fset.Position(pos))
			if uses, ok := tracker.Uses[varName]; ok {
				for _, usePos := range uses {
					fmt.Printf("    Used at: %v\n", fset.Position(usePos))
				}
			}
		}
		fmt.Println()
	}
}
