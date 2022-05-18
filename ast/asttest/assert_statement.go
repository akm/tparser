package asttest

import (
	"testing"

	"github.com/akm/tparser/ast"
	"github.com/stretchr/testify/assert"
)

func AssertCompoundStmt(t *testing.T, expected, actual *ast.CompoundStmt) {
	if !assert.Equal(t, expected.StmtList, actual.StmtList) {
		AssertStmtList(t, expected.StmtList, actual.StmtList)
	}
}

func AssertStmtList(t *testing.T, expected, actual ast.StmtList) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertStatement(t, exp, act)
		}
	}
}

func AssertStatement(t *testing.T, expected, actual *ast.Statement) {
	if !assert.Equal(t, expected.LabelId, actual.LabelId) {
		AssertLbelId(t, expected.LabelId, actual.LabelId)
	}
	if !assert.Equal(t, expected.Body, actual.Body) {
		AssertStatementBody(t, expected.Body, actual.Body)
	}
}

func AssertStatementBody(t *testing.T, expected, actual ast.StatementBody) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.CallStatement:
		AssertCallStatement(t, exp, actual.(*ast.CallStatement))
	case *ast.AssignStatement:
		AssertAssignStatement(t, exp, actual.(*ast.AssignStatement))
	case *ast.InheritedStatement:
		AssertInheritedStatement(t, exp, actual.(*ast.InheritedStatement))
	case *ast.GotoStatement:
		AssertGotoStatement(t, exp, actual.(*ast.GotoStatement))
	case *ast.CompoundStmt:
		AssertCompoundStmt(t, exp, actual.(*ast.CompoundStmt))
	case ast.ConditionalStmt:
		AssertConditionalStmt(t, exp, actual.(ast.ConditionalStmt))
	case ast.LoopStmt:
		AssertLoopStmt(t, exp, actual.(ast.LoopStmt))
	case *ast.WithStmt:
		AssertWithStmt(t, exp, actual.(*ast.WithStmt))
	case *ast.TryExceptStmt:
		AssertTryExceptStmt(t, exp, actual.(*ast.TryExceptStmt))
	case *ast.TryFinallyStmt:
		AssertTryFinallyStmt(t, exp, actual.(*ast.TryFinallyStmt))
	case *ast.RaiseStmt:
		AssertRaiseStmt(t, exp, actual.(*ast.RaiseStmt))
	case *ast.AssemblerStatement:
		AssertAssemblerStatement(t, exp, actual.(*ast.AssemblerStatement))
	default:
		assert.Fail(t, "unexpected statement type: %T", expected)
	}
}

func AssertCallStatement(t *testing.T, expected, actual *ast.CallStatement) {
	if !assert.Equal(t, expected.Designator, actual.Designator) {
		AssertDesignator(t, expected.Designator, actual.Designator)
	}
	if !assert.Equal(t, expected.ExprList, actual.ExprList) {
		AssertExprList(t, expected.ExprList, actual.ExprList)
	}
}

func AssertAssignStatement(t *testing.T, expected, actual *ast.AssignStatement) {
	if !assert.Equal(t, expected.Designator, actual.Designator) {
		AssertDesignator(t, expected.Designator, actual.Designator)
	}
	if !assert.Equal(t, expected.Expression, actual.Expression) {
		AssertExpression(t, expected.Expression, actual.Expression)
	}
}

func AssertInheritedStatement(t *testing.T, expected, actual *ast.InheritedStatement) {
	if !assert.Equal(t, expected.Ref, actual.Ref) {
		AssertDeclaration(t, expected.Ref, actual.Ref)
	}
}

func AssertGotoStatement(t *testing.T, expected, actual *ast.GotoStatement) {
	if !assert.Equal(t, expected.LabelId, actual.LabelId) {
		AssertLbelId(t, expected.LabelId, actual.LabelId)
	}
	if !assert.Equal(t, expected.Ref, actual.Ref) {
		AssertDeclaration(t, expected.Ref, actual.Ref)
	}
}

func AssertStructStmt(t *testing.T, expected, actual ast.StructStmt) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.CompoundStmt:
		AssertCompoundStmt(t, exp, actual.(*ast.CompoundStmt))
	case ast.ConditionalStmt:
		AssertConditionalStmt(t, exp, actual.(ast.ConditionalStmt))
	case ast.LoopStmt:
		AssertLoopStmt(t, exp, actual.(ast.LoopStmt))
	case *ast.WithStmt:
		AssertWithStmt(t, exp, actual.(*ast.WithStmt))
	case *ast.TryExceptStmt:
		AssertTryExceptStmt(t, exp, actual.(*ast.TryExceptStmt))
	case *ast.TryFinallyStmt:
		AssertTryFinallyStmt(t, exp, actual.(*ast.TryFinallyStmt))
	case *ast.RaiseStmt:
		AssertRaiseStmt(t, exp, actual.(*ast.RaiseStmt))
	case *ast.AssemblerStatement:
		AssertAssemblerStatement(t, exp, actual.(*ast.AssemblerStatement))
	default:
		assert.Fail(t, "unexpected statement type: %T", expected)
	}
}

func AssertConditionalStmt(t *testing.T, expected, actual ast.ConditionalStmt) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.IfStmt:
		AssertIfStmt(t, exp, actual.(*ast.IfStmt))
	case *ast.CaseStmt:
		AssertCaseStmt(t, exp, actual.(*ast.CaseStmt))
	default:
		assert.Fail(t, "unexpected conditional statement type: %T", expected)
	}
}

func AssertIfStmt(t *testing.T, expected, actual *ast.IfStmt) {
	if !assert.Equal(t, expected.Condition, actual.Condition) {
		AssertExpression(t, expected.Condition, actual.Condition)
	}
	if !assert.Equal(t, expected.Then, actual.Then) {
		AssertStatement(t, expected.Then, actual.Then)
	}
	if !assert.Equal(t, expected.Else, actual.Else) {
		AssertStatement(t, expected.Else, actual.Else)
	}
}

func AssertCaseStmt(t *testing.T, expected, actual *ast.CaseStmt) {
	if !assert.Equal(t, expected.Expression, actual.Expression) {
		AssertExpression(t, expected.Expression, actual.Expression)
	}
	if !assert.Equal(t, expected.Selectors, actual.Selectors) {
		AssertCaseSelectors(t, expected.Selectors, actual.Selectors)
	}
	if !assert.Equal(t, expected.Else, actual.Else) {
		AssertStmtList(t, expected.Else, actual.Else)
	}
}

func AssertCaseSelectors(t *testing.T, expected, actual ast.CaseSelectors) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i := range expected {
		exp := expected[i]
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertCaseSelector(t, exp, act)
		}
	}
}

func AssertCaseSelector(t *testing.T, expected, actual *ast.CaseSelector) {
	if !assert.Equal(t, expected.Labels, actual.Labels) {
		AssertCaseLabels(t, expected.Labels, actual.Labels)
	}
	if !assert.Equal(t, expected.Statement, actual.Statement) {
		AssertStatement(t, expected.Statement, actual.Statement)
	}
}

func AssertCaseLabels(t *testing.T, expected, actual ast.CaseLabels) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i, exp := range expected {
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertCaseLabel(t, exp, act)
		}
	}
}

func AssertCaseLabel(t *testing.T, expected, actual *ast.CaseLabel) {
	if !assert.Equal(t, expected.ConstExpr, actual.ConstExpr) {
		AssertConstExpr(t, expected.ConstExpr, actual.ConstExpr)
	}
	if !assert.Equal(t, expected.ExtraConstExpr, actual.ExtraConstExpr) {
		AssertConstExpr(t, expected.ExtraConstExpr, actual.ExtraConstExpr)
	}
}

func AssertLoopStmt(t *testing.T, expected, actual ast.LoopStmt) {
	if !assert.IsType(t, expected, actual) {
		return
	}
	switch exp := expected.(type) {
	case *ast.RepeatStmt:
		AssertRepeatStmt(t, exp, actual.(*ast.RepeatStmt))
	case *ast.WhileStmt:
		AssertWhileStmt(t, exp, actual.(*ast.WhileStmt))
	case *ast.ForStmt:
		AssertForStmt(t, exp, actual.(*ast.ForStmt))
	default:
		assert.Fail(t, "unexpected loop statement type: %T", expected)
	}
}

func AssertRepeatStmt(t *testing.T, expected, actual *ast.RepeatStmt) {
	if !assert.Equal(t, expected.StmtList, actual.StmtList) {
		AssertStmtList(t, expected.StmtList, actual.StmtList)
	}
	if !assert.Equal(t, expected.Condition, actual.Condition) {
		AssertExpression(t, expected.Condition, actual.Condition)
	}
}

func AssertWhileStmt(t *testing.T, expected, actual *ast.WhileStmt) {
	if !assert.Equal(t, expected.Condition, actual.Condition) {
		AssertExpression(t, expected.Condition, actual.Condition)
	}
	if !assert.Equal(t, expected.Statement, actual.Statement) {
		AssertStatement(t, expected.Statement, actual.Statement)
	}
}

func AssertForStmt(t *testing.T, expected, actual *ast.ForStmt) {
	if !assert.Equal(t, expected.QualId, actual.QualId) {
		AssertQualId(t, expected.QualId, actual.QualId)
	}
	if !assert.Equal(t, expected.Initial, actual.Initial) {
		AssertExpression(t, expected.Initial, actual.Initial)
	}
	if !assert.Equal(t, expected.Terminal, actual.Terminal) {
		AssertExpression(t, expected.Terminal, actual.Terminal)
	}
	assert.Equal(t, expected.Down, actual.Down)
	if !assert.Equal(t, expected.Statement, actual.Statement) {
		AssertStatement(t, expected.Statement, actual.Statement)
	}
}

func AssertWithStmt(t *testing.T, expected, actual *ast.WithStmt) {
	if !assert.Equal(t, expected.Objects, actual.Objects) {
		AssertQualIds(t, expected.Objects, actual.Objects)
	}
	if !assert.Equal(t, expected.Statement, actual.Statement) {
		AssertStatement(t, expected.Statement, actual.Statement)
	}
}

func AssertTryExceptStmt(t *testing.T, expected, actual *ast.TryExceptStmt) {
	if !assert.Equal(t, expected.Statements, actual.Statements) {
		AssertStmtList(t, expected.Statements, actual.Statements)
	}
	if !assert.Equal(t, expected.ExceptionBlock, actual.ExceptionBlock) {
		AssertExceptionBlock(t, expected.ExceptionBlock, actual.ExceptionBlock)
	}
}

func AssertExceptionBlock(t *testing.T, expected, actual *ast.ExceptionBlock) {
	if !assert.Equal(t, expected.Handlers, actual.Handlers) {
		AssertExceptionBlockHandlers(t, expected.Handlers, actual.Handlers)
	}
	if !assert.Equal(t, expected.Else, actual.Else) {
		AssertStmtList(t, expected.Else, actual.Else)
	}
}

func AssertExceptionBlockHandlers(t *testing.T, expected, actual ast.ExceptionBlockHandlers) {
	if !assert.Equal(t, len(expected), len(actual)) {
		return
	}
	for i := range expected {
		exp := expected[i]
		act := actual[i]
		if !assert.Equal(t, exp, act) {
			AssertExceptionHandler(t, exp, act)
		}
	}
}

func AssertExceptionHandler(t *testing.T, expected, actual *ast.ExceptionBlockHandler) {
	if !assert.Equal(t, expected.Decl, actual.Decl) {
		AssertExceptionBlockHandlerDecl(t, expected.Decl, actual.Decl)
	}
	if !assert.Equal(t, expected.Statement, actual.Statement) {
		AssertStatement(t, expected.Statement, actual.Statement)
	}
}

func AssertExceptionBlockHandlerDecl(t *testing.T, expected, actual *ast.ExceptionBlockHandlerDecl) {
	if !assert.Equal(t, expected.Ident, actual.Ident) {
		AssertIdent(t, expected.Ident, actual.Ident)
	}
	if !assert.Equal(t, expected.Type, actual.Type) {
		AssertType(t, expected.Type, actual.Type)
	}
}

func AssertTryFinallyStmt(t *testing.T, expected, actual *ast.TryFinallyStmt) {
	if !assert.Equal(t, expected.Statements1, actual.Statements1) {
		AssertStmtList(t, expected.Statements1, actual.Statements1)
	}
	if !assert.Equal(t, expected.Statements2, actual.Statements2) {
		AssertStmtList(t, expected.Statements2, actual.Statements2)
	}
}

func AssertRaiseStmt(t *testing.T, expected, actual *ast.RaiseStmt) {
	if !assert.Equal(t, expected.Object, actual.Object) {
		AssertExpression(t, expected.Object, actual.Object)
	}
	if !assert.Equal(t, expected.Address, actual.Address) {
		AssertExpression(t, expected.Address, actual.Address)
	}
}

func AssertAssemblerStatement(t *testing.T, expected, actual *ast.AssemblerStatement) {
	// Do nothing
}
