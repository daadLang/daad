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
		printLabelWithNode(childPrefix(prefix, isLast), testLast, "test:", n.Test)
		// Body
		bodyIsLast := len(n.Orelse) == 0
		printLabelWithNodes(childPrefix(prefix, isLast), bodyIsLast, "body:", nodesFromStmts(n.Body))
		// Orelse
		if len(n.Orelse) > 0 {
			printLabelWithNodes(childPrefix(prefix, isLast), true, "orelse:", nodesFromStmts(n.Orelse))
		}

	case *ForStmt:
		printHeader(prefix, isLast, "ForStmt")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "target:", n.Target)
		printLabelWithNode(base, false, "iter:", n.Iter)
		printLabelWithNodes(base, true, "body:", nodesFromStmts(n.Body))
		if len(n.Orelse) > 0 {
			printLabelWithNodes(base, true, "orelse:", nodesFromStmts(n.Orelse))
		}

	case *WhileStmt:
		printHeader(prefix, isLast, "WhileStmt")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "test:", n.Test)
		printLabelWithNodes(base, true, "body:", nodesFromStmts(n.Body))
		if len(n.Orelse) > 0 {
			printLabelWithNodes(base, true, "orelse:", nodesFromStmts(n.Orelse))
		}

	case *ExprStmt:
		printHeader(prefix, isLast, "ExprStmt")
		printLabelWithNode(childPrefix(prefix, isLast), true, "value:", n.Value)

	case *ReturnStmt:
		printHeader(prefix, isLast, "ReturnStmt")
		if n.Value != nil {
			printLabelWithNode(childPrefix(prefix, isLast), true, "value:", n.Value)
		}

	case *BreakStmt:
		printHeader(prefix, isLast, "BreakStmt")

	case *ContinueStmt:
		printHeader(prefix, isLast, "ContinueStmt")

	case *RepeatStmt:
		printHeader(prefix, isLast, "RepeatStmt")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "times:", n.Times)
		printLabelWithNodes(base, true, "body:", nodesFromStmts(n.Body))

	case *AugmentedAssignStmt:
		printHeader(prefix, isLast, fmt.Sprintf("AugmentedAssignStmt (op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "target:", n.Target)
		printLabelWithNode(base, true, "value:", n.Value)

	case *FunctionDefStmt:
		printHeader(prefix, isLast, fmt.Sprintf("FunctionDefStmt (name: %s, args: %v)", n.Name, n.Args))
		base := childPrefix(prefix, isLast)
		if len(n.Defaults) > 0 {
			printLabelWithNodes(base, false, "defaults:", nodesFromExprs(n.Defaults))
		}
		printLabelWithNodes(base, true, "body:", nodesFromStmts(n.Body))

	// Expressions
	case *Name:
		printHeader(prefix, isLast, fmt.Sprintf("Name(%s)", n.Id))

	case *Constant:
		printHeader(prefix, isLast, fmt.Sprintf("Constant(%v)", n.Value))

	case *Subscript:
		printHeader(prefix, isLast, "Subscript")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "value:", n.Value)
		printLabelWithNode(base, true, "index:", n.Index)

	case *Dict:
		printHeader(prefix, isLast, "Dict")
		base := childPrefix(prefix, isLast)
		for i := range n.Keys {
			keyLast := i == len(n.Keys)-1
			printLabelWithNode(base, !keyLast, "key:", n.Keys[i])
			printLabelWithNode(base, keyLast, "value:", n.Values[i])
		}

	case *List:
		printHeader(prefix, isLast, "List")
		printLabelWithNodes(childPrefix(prefix, isLast), true, "elements:", nodesFromExprs(n.Elements))

	case *Tuple:
		printHeader(prefix, isLast, "Tuple")
		printLabelWithNodes(childPrefix(prefix, isLast), true, "elements:", nodesFromExprs(n.Elements))

	case *Call:
		printHeader(prefix, isLast, "Call")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, len(n.Args) == 0, "Func:", n.Func)
		if len(n.Args) > 0 {
			printLabelWithNodes(base, true, "args:", nodesFromExprs(n.Args))
		}

	case *Compare:
		printHeader(prefix, isLast, fmt.Sprintf("Compare (op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "left:", n.Left)
		printLabelWithNode(base, true, "right:", n.Comparator)

	case *Assign:
		printHeader(prefix, isLast, "Assign")
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "target:", n.Target)
		printLabelWithNode(base, true, "value:", n.Value)

	case *BinOp:
		printHeader(prefix, isLast, fmt.Sprintf("BinOp (op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "left:", n.Left)
		printLabelWithNode(base, true, "right:", n.Right)

	case *UnaryOp:
		printHeader(prefix, isLast, fmt.Sprintf("UnaryOp (op: %s)", n.Op))
		printLabelWithNode(childPrefix(prefix, isLast), true, "expr:", n.Expr)

	case *BoolOp:
		printHeader(prefix, isLast, fmt.Sprintf("BoolOp (op: %s)", n.Op))
		base := childPrefix(prefix, isLast)
		printLabelWithNode(base, false, "left:", n.Left)
		printLabelWithNode(base, true, "right:", n.Right)

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
