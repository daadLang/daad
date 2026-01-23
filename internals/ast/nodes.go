package ast

import (
	"github.com/daadLang/daad/internals/lexer"
)

// TODO: read `https://docs.python.org/3/library/ast.html` for nodes

// ? ==========================================
// ? BASE
// ? ==========================================

type Node interface{}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

// module
type Module struct {
	Body []Stmt
}

// ? ==========================================
// ? STATEMENTS
// ? ==========================================
type IfStmt struct {
	Test   Expr
	Body   []Stmt
	Orelse []Stmt
}

func (*IfStmt) stmtNode() {}

type ForStmt struct {
	Target Expr
	Iter   Expr
	Body   []Stmt
	Orelse []Stmt
}

func (*ForStmt) stmtNode() {}

type WhileStmt struct {
	Test   Expr
	Body   []Stmt
	Orelse []Stmt
}

func (*WhileStmt) stmtNode() {}

// e.g a = 5
type ExprStmt struct {
	Value Expr
}

func (*ExprStmt) stmtNode() {}

type ReturnStmt struct {
	Value Expr
}

func (*ReturnStmt) stmtNode() {}

type RepeatStmt struct {
	Times Expr
	Body  []Stmt
}

func (*RepeatStmt) stmtNode() {}

type AugmentedAssignStmt struct {
	Target Expr
	Op     lexer.TokenType
	Value  Expr
}

func (*AugmentedAssignStmt) stmtNode() {}

type FunctionDefStmt struct {
	Name     string
	Args     []string
	Defaults []Expr
	Body     []Stmt
}

func (*FunctionDefStmt) stmtNode() {}

type AssignStmt struct {
	Target Name
	Value  Expr
}

func (*AssignStmt) stmtNode() {}

type BreakStmt struct{}

func (*BreakStmt) stmtNode() {}

type ContinueStmt struct{}

func (*ContinueStmt) stmtNode() {}

// ? ==========================================
// ? EXPRESSIONS
// ? ==========================================

type Name struct {
	Id string
}

func (*Name) exprNode() {}

// e.g
type Constant struct {
	Value any
}

func (*Constant) exprNode() {}

// arr[0] or dict["key"]  (value is dict/arr, inddex is 0/"key")
type Subscript struct {
	Value Expr
	Index Expr
}

func (*Subscript) exprNode() {}

type Dict struct {
	Keys   []Expr
	Values []Expr
}

func (*Dict) exprNode() {}

type List struct {
	Elements []Expr
}

func (*List) exprNode() {}

type Tuple struct {
	Elements []Expr
}

func (*Tuple) exprNode() {}

type Compare struct {
	Left       Expr
	Op         lexer.TokenType // ">", "<", "=="
	Comparator Expr
}

func (*Compare) exprNode() {}

type BinOp struct {
	Left  Expr
	Op    lexer.TokenType // "+", "-", "*", "/"
	Right Expr
}

func (*BinOp) exprNode() {}

type UnaryOp struct {
	Op   lexer.TokenType // "-", "not"
	Expr Expr
}

func (*UnaryOp) exprNode() {}

type BoolOp struct {
	Left  Expr
	Op    lexer.TokenType // "and", "or"
	Right Expr
}

func (*BoolOp) exprNode() {}

type Kwarg struct {
	Name  string
	Value Expr
}

// Call represents a function call expression (e.g., func(args))
type Call struct {
	Func   Expr
	Args   []Expr  // positional args (e.g., func(1, 2))
	Kwargs []Kwarg // keyword args (name=value)
}

func (*Call) exprNode() {}

// Assign represents an assignment expression (e.g., a = 5)
type Assign struct {
	Target Expr
	Value  Expr
}

func (*Assign) exprNode() {}
