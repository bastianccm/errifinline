package errifinline

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func NewAnalyzer() (*analysis.Analyzer, error) {
	a, err := newAnalyzer()
	if err != nil {
		return nil, err
	}

	return &analysis.Analyzer{
		Name:     "errifinline",
		Doc:      "check for err assignments not inlined in if",
		Run:      a.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}, nil
}

type analyzer struct{}

func newAnalyzer() (*analyzer, error) {
	a := &analyzer{}

	return a, nil
}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	inspectorInfo := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.IfStmt)(nil), (*ast.AssignStmt)(nil)}

	inspectorInfo.Preorder(nodeFilter, a.AsCheckVisitor(pass))

	return nil, nil
}

func isSingleErrAssignment(info *types.Info, assign *ast.AssignStmt) (token.Pos, bool) {
	var pos token.Pos
	var singleErr = true
	for _, expr := range assign.Lhs {
		ident, ok := expr.(*ast.Ident)
		if !ok {
			continue
		}
		if ident.Name == "_" {
			continue
		}
		if info.TypeOf(ident) != nil && info.TypeOf(ident).String() != "error" {
			singleErr = false
			continue
		}
		if info.ObjectOf(ident) != nil {
			pos = info.ObjectOf(ident).Pos()
		}
	}
	return pos, singleErr
}

func (a *analyzer) AsCheckVisitor(pass *analysis.Pass) func(ast.Node) {
	singleErr := make(map[token.Pos]struct{})

	return func(n ast.Node) {
		assign, ok := n.(*ast.AssignStmt)
		if ok {
			ident, ok := isSingleErrAssignment(pass.TypesInfo, assign)
			if ok {
				singleErr[ident] = struct{}{}
			} else {
				delete(singleErr, ident)
			}
			return
		}

		ifstmt, ok := n.(*ast.IfStmt)
		if !ok {
			return
		}

		if ifstmt.Init != nil {
			return
		}

		cond, ok := ifstmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return
		}

		if pass.TypesInfo.TypeOf(cond.X).String() != "error" {
			return
		}

		ident, ok := cond.X.(*ast.Ident)
		if !ok {
			return
		}

		_, singleErr := singleErr[pass.TypesInfo.ObjectOf(ident).Pos()]
		if !singleErr {
			return
		}

		d := analysis.Diagnostic{
			Pos:      n.Pos(),
			End:      n.End(),
			Message:  "inline err assignment in if initializer",
			Category: "errifinline",
		}
		pass.Report(d)
	}
}
