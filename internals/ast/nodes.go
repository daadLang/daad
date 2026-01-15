package ast

import (
	"fmt"
	"strings"

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

type BreakStmt struct{}

func (*BreakStmt) stmtNode() {}

type ContinueStmt struct{}

func (*ContinueStmt) stmtNode() {}

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

// ? ==========================================
// ? AST PRINTING
// ? ==========================================

// PrintAST prints the AST tree in a readable format
func PrintAST(module Module) {
	fmt.Println("Module")
	for _, stmt := range module.Body {
		printNode(stmt, 1)
	}
}

func printNode(node Node, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch n := node.(type) {
	// Statements
	case *IfStmt:
		fmt.Printf("%sIfStmt\n", prefix)
		fmt.Printf("%s  Test:\n", prefix)
		printNode(n.Test, indent+2)
		fmt.Printf("%s  Body:\n", prefix)
		for _, stmt := range n.Body {
			printNode(stmt, indent+2)
		}
		if len(n.Orelse) > 0 {
			fmt.Printf("%s  Orelse:\n", prefix)
			for _, stmt := range n.Orelse {
				printNode(stmt, indent+2)
			}
		}

	case *ForStmt:
		fmt.Printf("%sForStmt\n", prefix)
		fmt.Printf("%s  Target:\n", prefix)
		printNode(n.Target, indent+2)
		fmt.Printf("%s  Iter:\n", prefix)
		printNode(n.Iter, indent+2)
		fmt.Printf("%s  Body:\n", prefix)
		for _, stmt := range n.Body {
			printNode(stmt, indent+2)
		}
		if len(n.Orelse) > 0 {
			fmt.Printf("%s  Orelse:\n", prefix)
			for _, stmt := range n.Orelse {
				printNode(stmt, indent+2)
			}
		}

	case *WhileStmt:
		fmt.Printf("%sWhileStmt\n", prefix)
		fmt.Printf("%s  Test:\n", prefix)
		printNode(n.Test, indent+2)
		fmt.Printf("%s  Body:\n", prefix)
		for _, stmt := range n.Body {
			printNode(stmt, indent+2)
		}
		if len(n.Orelse) > 0 {
			fmt.Printf("%s  Orelse:\n", prefix)
			for _, stmt := range n.Orelse {
				printNode(stmt, indent+2)
			}
		}

	case *ExprStmt:
		fmt.Printf("%sExprStmt\n", prefix)
		printNode(n.Value, indent+1)

	case *ReturnStmt:
		fmt.Printf("%sReturnStmt\n", prefix)
		if n.Value != nil {
			printNode(n.Value, indent+1)
		}

	case *BreakStmt:
		fmt.Printf("%sBreakStmt\n", prefix)

	case *ContinueStmt:
		fmt.Printf("%sContinueStmt\n", prefix)

	case *RepeatStmt:
		fmt.Printf("%sRepeatStmt\n", prefix)
		fmt.Printf("%s  Times:\n", prefix)
		printNode(n.Times, indent+2)
		fmt.Printf("%s  Body:\n", prefix)
		for _, stmt := range n.Body {
			printNode(stmt, indent+2)
		}

	case *AugmentedAssignStmt:
		fmt.Printf("%sAugmentedAssignStmt (Op: %s)\n", prefix, n.Op)
		fmt.Printf("%s  Target:\n", prefix)
		printNode(n.Target, indent+2)
		fmt.Printf("%s  Value:\n", prefix)
		printNode(n.Value, indent+2)

	case *FunctionDefStmt:
		fmt.Printf("%sFunctionDefStmt (Name: %s, Args: %v)\n", prefix, n.Name, n.Args)
		if len(n.Defaults) > 0 {
			fmt.Printf("%s  Defaults:\n", prefix)
			for _, def := range n.Defaults {
				printNode(def, indent+2)
			}
		}
		fmt.Printf("%s  Body:\n", prefix)
		for _, stmt := range n.Body {
			printNode(stmt, indent+2)
		}

	// Expressions
	case *Name:
		fmt.Printf("%sName(%s)\n", prefix, n.Id)

	case *Constant:
		fmt.Printf("%sConstant(%v)\n", prefix, n.Value)

	case *Subscript:
		fmt.Printf("%sSubscript\n", prefix)
		fmt.Printf("%s  Value:\n", prefix)
		printNode(n.Value, indent+2)
		fmt.Printf("%s  Index:\n", prefix)
		printNode(n.Index, indent+2)

	case *Dict:
		fmt.Printf("%sDict\n", prefix)
		for i := range n.Keys {
			fmt.Printf("%s  Key:\n", prefix)
			printNode(n.Keys[i], indent+2)
			fmt.Printf("%s  Value:\n", prefix)
			printNode(n.Values[i], indent+2)
		}

	case *List:
		fmt.Printf("%sList\n", prefix)
		for _, elem := range n.Elements {
			printNode(elem, indent+1)
		}

	case *Tuple:
		fmt.Printf("%sTuple\n", prefix)
		for _, elem := range n.Elements {
			printNode(elem, indent+1)
		}

	case *Call:
		fmt.Printf("%sCall\n", prefix)
		fmt.Printf("%s  Func:\n", prefix)
		printNode(n.Func, indent+2)
		if len(n.Args) > 0 {
			fmt.Printf("%s  Args:\n", prefix)
			for _, arg := range n.Args {
				printNode(arg, indent+2)
			}
		}

	case *Compare:
		fmt.Printf("%sCompare (Op: %s)\n", prefix, n.Op)
		fmt.Printf("%s  Left:\n", prefix)
		printNode(n.Left, indent+2)
		fmt.Printf("%s  Right:\n", prefix)
		printNode(n.Comparator, indent+2)

	case *Assign:
		fmt.Printf("%sAssign\n", prefix)
		fmt.Printf("%s  Target:\n", prefix)
		printNode(n.Target, indent+2)
		fmt.Printf("%s  Value:\n", prefix)
		printNode(n.Value, indent+2)

	case *BinOp:
		fmt.Printf("%sBinOp (Op: %s)\n", prefix, n.Op)
		fmt.Printf("%s  Left:\n", prefix)
		printNode(n.Left, indent+2)
		fmt.Printf("%s  Right:\n", prefix)
		printNode(n.Right, indent+2)

	case *UnaryOp:
		fmt.Printf("%sUnaryOp (Op: %s)\n", prefix, n.Op)
		printNode(n.Expr, indent+1)

	case *BoolOp:
		fmt.Printf("%sBoolOp (Op: %s)\n", prefix, n.Op)
		fmt.Printf("%s  Left:\n", prefix)
		printNode(n.Left, indent+2)
		fmt.Printf("%s  Right:\n", prefix)
		printNode(n.Right, indent+2)

	case *CallExpr:
		fmt.Printf("%sCallExpr\n", prefix)
		fmt.Printf("%s  Func:\n", prefix)
		printNode(n.Func, indent+2)
		if len(n.Args) > 0 {
			fmt.Printf("%s  Args:\n", prefix)
			for _, arg := range n.Args {
				printNode(arg, indent+2)
			}
		}

	default:
		fmt.Printf("%s<unknown node: %T>\n", prefix, node)
	}
}
