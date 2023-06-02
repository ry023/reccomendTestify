package recommendTestify

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = `recommendTestify is a linter that outputs a message recommending the use of Testify when testing.T is used directly.
This is a unique rule in hosborn development.`

var Analyzer = &analysis.Analyzer{
	Name: "recommendTestify",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	for _, method := range findMethods(pass) {
		if method.ReceiverType.String() == "*testing.T" {
			switch method.MethodIdent.Name {
			case "Errorf":
				pass.Reportf(method.CallExpr.Pos(), "We recommend using Testify (github.com/stretchr/testify). Use Testify's assertion methods or wrap by `assert.Failf` method.")
			case "Error":
				pass.Reportf(method.CallExpr.Pos(), "We recommend using Testify (github.com/stretchr/testify). Use Testify's assertion methods or wrap by `assert.Fail` method.")
			case "Fatal":
				pass.Reportf(method.CallExpr.Pos(), "We recommend using Testify (github.com/stretchr/testify). Use Testify's assertion methods or wrap by `assert.FailNow` method.")
			case "Fatalf":
				pass.Reportf(method.CallExpr.Pos(), "We recommend using Testify (github.com/stretchr/testify). Use Testify's assertion methods or wrap by `assert.FailNowf` method.")
			}
		}
	}
	return nil, nil
}

type MethodCall struct {
	CallExpr     *ast.CallExpr
	ReceiverExpr ast.Expr
	ReceiverType types.Type
	MethodIdent  *ast.Ident
}

func findMethods(pass *analysis.Pass) []*MethodCall {
	methods := []*MethodCall{}

	// Find functions or methods calling
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Check if the node is method
		methodExpr, isMethod := callExpr.Fun.(*ast.SelectorExpr)
		if !isMethod {
			return
		}

		receiverExpr := methodExpr.X
		methodIdent := methodExpr.Sel

		methods = append(
			methods,
			&MethodCall{
				CallExpr:     callExpr,
				ReceiverExpr: receiverExpr,
				ReceiverType: pass.TypesInfo.TypeOf(receiverExpr),
				MethodIdent:  methodIdent,
			},
		)
	})
	return methods
}
