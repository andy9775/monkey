package evaluator

import (
	"github.com/andy9775/interpreter/ast"
	"github.com/andy9775/interpreter/object"
)

// Eval takes in an AST node, determines it's type and returns the
// resulting object representation of that type
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program: // evaluate the statements
		return evalStatements(node.Statements)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}
