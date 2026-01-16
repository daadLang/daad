package ast

import (
	"fmt"
	"strings"
)

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

	default:
		fmt.Printf("%s<unknown node: %T>\n", prefix, node)
	}
}
