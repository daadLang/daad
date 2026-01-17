package ast

import (
	"fmt"
)

// ? ==========================================
// ? AST PRINTING
// ? ==========================================

func PrintAST(module Module) {
	fmt.Println("Module")
	for i, stmt := range module.Body {
		last := i == len(module.Body)-1
		printNode(stmt, "", last)
	}
}

func printHeader(prefix string, isLast bool, header string) {
	conn := "├── "
	if isLast {
		conn = "└── "
	}
	fmt.Printf("%s%s%s\n", prefix, conn, header)
}

func childPrefix(prefix string, isLast bool) string {
	if isLast {
		return prefix + "    "
	}
	return prefix + "│   "
}

func printLabelWithNode(prefix string, fieldIsLast bool, label string, node Node) {
	printHeader(prefix, fieldIsLast, label)
	newPref := childPrefix(prefix, fieldIsLast)
	// single child under the label -> it's the last sibling inside that label
	printNode(node, newPref, true)
}

func printLabelWithNodes(prefix string, fieldIsLast bool, label string, nodes []Node) {
	printHeader(prefix, fieldIsLast, label)
	newPref := childPrefix(prefix, fieldIsLast)
	for i, n := range nodes {
		last := i == len(nodes)-1
		printNode(n, newPref, last)
	}
}

func printNode(node Node, prefix string, isLast bool) {
	switch n := node.(type) {
	// Statements
	case *IfStmt:
		printHeader(prefix, isLast, "IfStmt")
		// children: Test, Body, (optional) Orelse
		// Test
		testLast := len(n.Body) == 0 && len(n.Orelse) == 0
		printLabelWithNode(childPrefix(prefix, isLast), testLast, "Test:", n.Test)
		// Body
		bodyIsLast := len(n.Orelse) == 0
		printLabelWithNodes(childPrefix(prefix, isLast), bodyIsLast, "Body:", nodesFromStmts(n.Body))
		// Orelse
		if len(n.Orelse) > 0 {
			printLabelWithNodes(childPrefix(prefix, isLast), true, "Orelse:", nodesFromStmts(n.Orelse))
		}

	case *ForStmt:
		printHeader(prefix, isLast, "ForStmt")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Target:", n.Target)
		printLabelWithNode(base, false, "Iter:", n.Iter)
		printLabelWithNodes(base, true, "Body:", nodesFromStmts(n.Body))
		if len(n.Orelse) > 0 {
			printLabelWithNodes(base, true, "Orelse:", nodesFromStmts(n.Orelse))
		}

	case *WhileStmt:
		printHeader(prefix, isLast, "WhileStmt")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Test:", n.Test)
		printLabelWithNodes(base, true, "Body:", nodesFromStmts(n.Body))
		if len(n.Orelse) > 0 {
			printLabelWithNodes(base, true, "Orelse:", nodesFromStmts(n.Orelse))
		}

	case *ExprStmt:
		printHeader(prefix, isLast, "ExprStmt")
		printLabelWithNode(childPrefix(prefix, isLast), true, "Value:", n.Value)

	case *ReturnStmt:
		printHeader(prefix, isLast, "ReturnStmt")
		if n.Value != nil {
			printLabelWithNode(childPrefix(prefix, isLast), true, "Value:", n.Value)
		}

	case *BreakStmt:
		printHeader(prefix, isLast, "BreakStmt")

	case *ContinueStmt:
		printHeader(prefix, isLast, "ContinueStmt")

	case *RepeatStmt:
		printHeader(prefix, isLast, "RepeatStmt")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Times:", n.Times)
		printLabelWithNodes(base, true, "Body:", nodesFromStmts(n.Body))

	case *AugmentedAssignStmt:
		printHeader(prefix, isLast, fmt.Sprintf("AugmentedAssignStmt (Op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Target:", n.Target)
		printLabelWithNode(base, true, "Value:", n.Value)

	case *FunctionDefStmt:
		printHeader(prefix, isLast, fmt.Sprintf("FunctionDefStmt (Name: %s, Args: %v)", n.Name, n.Args))
		base := childPrefix(prefix, isLast)
		if len(n.Defaults) > 0 {
			printLabelWithNodes(base, false, "Defaults:", nodesFromExprs(n.Defaults))
		}
		printLabelWithNodes(base, true, "Body:", nodesFromStmts(n.Body))

	// Expressions
	case *Name:
		printHeader(prefix, isLast, fmt.Sprintf("Name(%s)", n.Id))

	case *Constant:
		printHeader(prefix, isLast, fmt.Sprintf("Constant(%v)", n.Value))

	case *Subscript:
		printHeader(prefix, isLast, "Subscript")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Value:", n.Value)
		printLabelWithNode(base, true, "Index:", n.Index)

	case *Dict:
		printHeader(prefix, isLast, "Dict")
		base := childPrefix(prefix, isLast)
		for i := range n.Keys {
			keyLast := i == len(n.Keys)-1
			printLabelWithNode(base, !keyLast, "Key:", n.Keys[i])
			printLabelWithNode(base, keyLast, "Value:", n.Values[i])
		}

	case *List:
		printHeader(prefix, isLast, "List")
		printLabelWithNodes(childPrefix(prefix, isLast), true, "Elements:", nodesFromExprs(n.Elements))

	case *Tuple:
		printHeader(prefix, isLast, "Tuple")
		printLabelWithNodes(childPrefix(prefix, isLast), true, "Elements:", nodesFromExprs(n.Elements))

	case *Call:
		printHeader(prefix, isLast, "Call")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, len(n.Args) == 0, "Func:", n.Func)
		if len(n.Args) > 0 {
			printLabelWithNodes(base, true, "Args:", nodesFromExprs(n.Args))
		}

	case *Compare:
		printHeader(prefix, isLast, fmt.Sprintf("Compare (Op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Left:", n.Left)
		printLabelWithNode(base, true, "Right:", n.Comparator)

	case *Assign:
		printHeader(prefix, isLast, "Assign")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Target:", n.Target)
		printLabelWithNode(base, true, "Value:", n.Value)

	case *BinOp:
		printHeader(prefix, isLast, fmt.Sprintf("BinOp (Op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Left:", n.Left)
		printLabelWithNode(base, true, "Right:", n.Right)

	case *UnaryOp:
		printHeader(prefix, isLast, fmt.Sprintf("UnaryOp (Op: %s)", n.Op))
		printLabelWithNode(childPrefix(prefix, isLast), true, "Expr:", n.Expr)

	case *BoolOp:
		printHeader(prefix, isLast, fmt.Sprintf("BoolOp (Op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "Left:", n.Left)
		printLabelWithNode(base, true, "Right:", n.Right)

	default:
		printHeader(prefix, isLast, fmt.Sprintf("<unknown node: %T>", node))
	}
}

// helpers to convert specific typed slices into []Node for printing
func nodesFromStmts(stmts []Stmt) []Node {
	nodes := make([]Node, len(stmts))
	for i := range stmts {
		nodes[i] = stmts[i]
	}
	return nodes
}

func nodesFromExprs(exprs []Expr) []Node {
	nodes := make([]Node, len(exprs))
	for i := range exprs {
		nodes[i] = exprs[i]
	}
	return nodes
}
