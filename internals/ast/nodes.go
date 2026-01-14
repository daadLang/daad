package ast

import "github.com/daadLang/daad/internals/lexer"

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

type FunctionDefStmt struct {
	Name     string
	Args     []string
	Defaults []Expr
	Body     []Stmt
}

func (*FunctionDefStmt) stmtNode() {}

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

type Call struct {
	Func Expr
	Args []Expr
}

func (*Call) exprNode() {}

type Compare struct {
	Left       Expr
	Op         lexer.TokenType // ">", "<", "=="
	Comparator Expr
}

func (*Compare) exprNode() {}

type Assign struct {
	Target Expr
	Value  Expr
}

func (*Assign) exprNode() {}

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

type CallExpr struct {
	Func Expr
	Args []Expr
}

func (*CallExpr) exprNode() {}
