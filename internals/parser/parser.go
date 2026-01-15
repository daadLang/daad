package parser

import (
	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
)

type Parser struct {
	Tokens []lexer.Token
	Pos    int
	Module ast.Module
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		Tokens: tokens,
		Pos:    0,
		Module: ast.Module{Body: []ast.Stmt{}},
	}
}

func (p *Parser) peek() lexer.Token {
	if p.Pos >= len(p.Tokens) {
		return lexer.Token{Type: lexer.EOF, Value: ""}
	}
	return p.Tokens[p.Pos]
}

func (p *Parser) advance() lexer.Token {
	token := p.peek()
	p.Pos++
	return token
}

func (p *Parser) expect(tokenType lexer.TokenType) bool {
	token := p.peek()
	if token.Type != tokenType {
		return false
	}
	return true
}

func (p *Parser) Parse() ast.Module {
	for p.Pos < len(p.Tokens) {
		stmt := p.parseStatement()
		if stmt != nil {
			p.Module.Body = append(p.Module.Body, stmt)
		}
	}
	return p.Module
}

func (p *Parser) parseStatement() ast.Stmt {
	token := p.peek()

	switch token.Type {
	// Compound statements
	case lexer.IF:
		return p.parseIfStmt()
	case lexer.WHILE:
		return p.parseWhileStmt()
	case lexer.FOR:
		return p.parseForStmt()
	case lexer.REPEAT:
		return p.parseRepeatStmt()
	case lexer.FUNC:
		return p.parseFunctionDef()

	// Simple statements
	case lexer.RETURN:
		return p.parseReturnStmt()
	case lexer.BREAK:
		return p.parseBreakStmt()
	case lexer.CONTINUE:
		return p.parseContinueStmt()

	case lexer.NEWLINE:
		p.advance() // skip empty lines
		return nil
	case lexer.EOF:
		p.advance()
		return nil
	case lexer.COMMENT:
		p.advance() // skip comments
		return nil

	default:
		// Could be: expr_stmt, assignment_stmt, or augmented_assign_stmt
		return p.parseExprOrAssignStmt()
	}
}

// parseExprOrAssignStmt handles expression statements, assignments, and augmented assignments
func (p *Parser) parseExprOrAssignStmt() ast.Stmt {
	expr := p.parseExpression()

	token := p.peek()
	switch token.Type {
	case lexer.ASSIGN:
		p.advance()
		value := p.parseExpression()
		p.consumeNewline()
		return &ast.ExprStmt{Value: &ast.Assign{Target: expr, Value: value}}

	case lexer.PLUS_ASSIGN, lexer.MINUS_ASSIGN, lexer.MULT_ASSIGN,
		lexer.DIVIDE_ASSIGN, lexer.MOD_ASSIGN, lexer.POWER_ASSIGN:
		op := p.advance()
		value := p.parseExpression()
		p.consumeNewline()
		return &ast.AugmentedAssignStmt{Target: expr, Op: op.Type, Value: value}

	default:
		p.consumeNewline()
		return &ast.ExprStmt{Value: expr}
	}
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	p.advance() // consume RETURN

	var value ast.Expr = nil
	if !p.isAtEnd() && p.peek().Type != lexer.NEWLINE {
		value = p.parseExpression()
	}
	p.consumeNewline()
	return &ast.ReturnStmt{Value: value}
}

func (p *Parser) parseBreakStmt() *ast.BreakStmt {
	p.advance() // consume BREAK
	p.consumeNewline()
	return &ast.BreakStmt{}
}

func (p *Parser) parseContinueStmt() *ast.ContinueStmt {
	p.advance() // consume CONTINUE
	p.consumeNewline()
	return &ast.ContinueStmt{}
}

func (p *Parser) parseIfStmt() *ast.IfStmt {
	p.advance() // consume IF

	test := p.parseExpression()
	p.expectAndAdvance(lexer.COLON)
	body := p.parseSuite()

	var orelse []ast.Stmt

	// Handle elif chains
	for p.peek().Type == lexer.ELIF {
		p.advance() // consume ELIF
		elifTest := p.parseExpression()
		p.expectAndAdvance(lexer.COLON)
		elifBody := p.parseSuite()
		orelse = []ast.Stmt{&ast.IfStmt{Test: elifTest, Body: elifBody, Orelse: nil}}
	}

	// Handle else
	if p.peek().Type == lexer.ELSE {
		p.advance() // consume ELSE
		p.expectAndAdvance(lexer.COLON)
		elseBody := p.parseSuite()
		if len(orelse) > 0 {
			// Attach else to the last elif
			lastElif := orelse[0].(*ast.IfStmt)
			for lastElif.Orelse != nil && len(lastElif.Orelse) > 0 {
				if elif, ok := lastElif.Orelse[0].(*ast.IfStmt); ok {
					lastElif = elif
				} else {
					break
				}
			}
			lastElif.Orelse = elseBody
		} else {
			orelse = elseBody
		}
	}

	return &ast.IfStmt{Test: test, Body: body, Orelse: orelse}
}

func (p *Parser) parseWhileStmt() *ast.WhileStmt {
	p.advance() // consume WHILE

	test := p.parseExpression()
	p.expectAndAdvance(lexer.COLON)
	body := p.parseSuite()

	return &ast.WhileStmt{Test: test, Body: body, Orelse: nil}
}

func (p *Parser) parseForStmt() *ast.ForStmt {
	p.advance() // consume FOR

	targetToken := p.advance() // get the loop variable NAME
	target := &ast.Name{Id: targetToken.Value}

	p.expectAndAdvance(lexer.IN)
	iter := p.parseExpression()
	p.expectAndAdvance(lexer.COLON)
	body := p.parseSuite()

	return &ast.ForStmt{Target: target, Iter: iter, Body: body, Orelse: nil}
}

func (p *Parser) parseRepeatStmt() *ast.RepeatStmt {
	p.advance() // consume REPEAT

	times := p.parseExpression()
	p.expectAndAdvance(lexer.TIMES)
	p.expectAndAdvance(lexer.COLON)
	body := p.parseSuite()

	return &ast.RepeatStmt{Times: times, Body: body}
}

func (p *Parser) parseFunctionDef() *ast.FunctionDefStmt {
	p.advance() // consume FUNC

	nameToken := p.advance()
	name := nameToken.Value

	p.expectAndAdvance(lexer.LPAREN)

	var args []string
	var defaults []ast.Expr

	if p.peek().Type != lexer.RPAREN {
		for {
			paramToken := p.advance()
			args = append(args, paramToken.Value)

			// Check for default value
			if p.peek().Type == lexer.ASSIGN {
				p.advance()
				defaultVal := p.parseExpression()
				defaults = append(defaults, defaultVal)
			}

			if p.peek().Type != lexer.COMMA {
				break
			}
			p.advance() // consume COMMA
		}
	}

	p.expectAndAdvance(lexer.RPAREN)
	p.expectAndAdvance(lexer.COLON)
	body := p.parseSuite()

	return &ast.FunctionDefStmt{Name: name, Args: args, Defaults: defaults, Body: body}
}

func (p *Parser) parseSuite() []ast.Stmt {
	var stmts []ast.Stmt

	// Skip newline before INDENT
	if p.peek().Type == lexer.NEWLINE {
		p.advance()
	}

	// Expect INDENT
	if p.peek().Type == lexer.INDENT {
		p.advance()

		for p.peek().Type != lexer.DEDENT && p.peek().Type != lexer.EOF {
			stmt := p.parseStatement()
			if stmt != nil {
				stmts = append(stmts, stmt)
			}
		}

		// Consume DEDENT
		if p.peek().Type == lexer.DEDENT {
			p.advance()
		}
	} else {
		// Single-line suite (simple_stmt NEWLINE)
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	return stmts
}

// Helper functions

func (p *Parser) expectAndAdvance(tokenType lexer.TokenType) lexer.Token {
	if p.peek().Type != tokenType {
		// For now, just advance anyway - could add error handling
	}
	return p.advance()
}

func (p *Parser) consumeNewline() {
	if p.peek().Type == lexer.NEWLINE {
		p.advance()
	}
}

func (p *Parser) isAtEnd() bool {
	return p.Pos >= len(p.Tokens) || p.peek().Type == lexer.EOF
}

// parseExpression - placeholder for expression parsing
// This should be implemented to handle the full expression grammar
func (p *Parser) parseExpression() ast.Expr {
	return p.parseOrExpr()
}

func (p *Parser) parseOrExpr() ast.Expr {
	left := p.parseAndExpr()

	for p.peek().Type == lexer.OR {
		op := p.advance()
		right := p.parseAndExpr()
		left = &ast.BoolOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseAndExpr() ast.Expr {
	left := p.parseNotExpr()

	for p.peek().Type == lexer.AND {
		op := p.advance()
		right := p.parseNotExpr()
		left = &ast.BoolOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseNotExpr() ast.Expr {
	if p.peek().Type == lexer.NOT {
		op := p.advance()
		expr := p.parseNotExpr()
		return &ast.UnaryOp{Op: op.Type, Expr: expr}
	}
	return p.parseComparison()
}

func (p *Parser) parseComparison() ast.Expr {
	left := p.parseBitorExpr()

	for p.isComparisonOp(p.peek().Type) {
		op := p.advance()
		right := p.parseBitorExpr()
		left = &ast.Compare{Left: left, Op: op.Type, Comparator: right}
	}

	return left
}

func (p *Parser) isComparisonOp(t lexer.TokenType) bool {
	return t == lexer.EQ || t == lexer.NEQ || t == lexer.LESS ||
		t == lexer.GREATER || t == lexer.LEQ || t == lexer.GEQ || t == lexer.IN
}

func (p *Parser) parseBitorExpr() ast.Expr {
	left := p.parseBitxorExpr()

	for p.peek().Type == lexer.BITWISE_OR {
		op := p.advance()
		right := p.parseBitxorExpr()
		left = &ast.BinOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseBitxorExpr() ast.Expr {
	left := p.parseBitandExpr()

	for p.peek().Type == lexer.BITWISE_XOR {
		op := p.advance()
		right := p.parseBitandExpr()
		left = &ast.BinOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseBitandExpr() ast.Expr {
	left := p.parseShiftExpr()

	for p.peek().Type == lexer.BITWISE_AND {
		op := p.advance()
		right := p.parseShiftExpr()
		left = &ast.BinOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseShiftExpr() ast.Expr {
	left := p.parseArithExpr()

	for p.peek().Type == lexer.LSHIFT || p.peek().Type == lexer.RSHIFT {
		op := p.advance()
		right := p.parseArithExpr()
		left = &ast.BinOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseArithExpr() ast.Expr {
	left := p.parseTerm()

	for p.peek().Type == lexer.PLUS || p.peek().Type == lexer.MINUS {
		op := p.advance()
		right := p.parseTerm()
		left = &ast.BinOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseTerm() ast.Expr {
	left := p.parseFactor()

	for p.peek().Type == lexer.MULT || p.peek().Type == lexer.DIVIDE ||
		p.peek().Type == lexer.FLOORDIV || p.peek().Type == lexer.MOD {
		op := p.advance()
		right := p.parseFactor()
		left = &ast.BinOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parseFactor() ast.Expr {
	token := p.peek()

	if token.Type == lexer.PLUS || token.Type == lexer.MINUS || token.Type == lexer.BITWISE_NOT {
		op := p.advance()
		expr := p.parseFactor()
		return &ast.UnaryOp{Op: op.Type, Expr: expr}
	}

	return p.parsePower()
}

func (p *Parser) parsePower() ast.Expr {
	left := p.parsePrimary()

	if p.peek().Type == lexer.POWER {
		op := p.advance()
		right := p.parseFactor() // right-associative
		return &ast.BinOp{Left: left, Op: op.Type, Right: right}
	}

	return left
}

func (p *Parser) parsePrimary() ast.Expr {
	atom := p.parseAtom()

	for {
		switch p.peek().Type {
		case lexer.LPAREN:
			// Function call
			p.advance()
			var args []ast.Expr
			if p.peek().Type != lexer.RPAREN {
				args = append(args, p.parseExpression())
				for p.peek().Type == lexer.COMMA {
					p.advance()
					args = append(args, p.parseExpression())
				}
			}
			p.expectAndAdvance(lexer.RPAREN)
			atom = &ast.Call{Func: atom, Args: args}

		case lexer.LBRACKET:
			// Subscript
			p.advance()
			index := p.parseExpression()
			p.expectAndAdvance(lexer.RBRACKET)
			atom = &ast.Subscript{Value: atom, Index: index}

		case lexer.DOT:
			// Attribute access
			p.advance()
			attrToken := p.advance()
			atom = &ast.Subscript{Value: atom, Index: &ast.Constant{Value: attrToken.Value}}

		default:
			return atom
		}
	}
}

func (p *Parser) parseAtom() ast.Expr {
	token := p.peek()

	switch token.Type {
	case lexer.NAME:
		p.advance()
		return &ast.Name{Id: token.Value}

	case lexer.NUMBER:
		p.advance()
		return &ast.Constant{Value: token.Value}

	case lexer.STRING:
		p.advance()
		return &ast.Constant{Value: token.Value}

	case lexer.TRUE:
		p.advance()
		return &ast.Constant{Value: true}

	case lexer.FALSE:
		p.advance()
		return &ast.Constant{Value: false}

	case lexer.LPAREN:
		p.advance()
		// Could be tuple or parenthesized expression
		if p.peek().Type == lexer.RPAREN {
			p.advance()
			return &ast.Tuple{Elements: []ast.Expr{}}
		}
		expr := p.parseExpression()
		if p.peek().Type == lexer.COMMA {
			// It's a tuple
			elements := []ast.Expr{expr}
			for p.peek().Type == lexer.COMMA {
				p.advance()
				if p.peek().Type == lexer.RPAREN {
					break
				}
				elements = append(elements, p.parseExpression())
			}
			p.expectAndAdvance(lexer.RPAREN)
			return &ast.Tuple{Elements: elements}
		}
		p.expectAndAdvance(lexer.RPAREN)
		return expr

	case lexer.LBRACKET:
		// List literal
		p.advance()
		var elements []ast.Expr
		if p.peek().Type != lexer.RBRACKET {
			elements = append(elements, p.parseExpression())
			for p.peek().Type == lexer.COMMA {
				p.advance()
				if p.peek().Type == lexer.RBRACKET {
					break
				}
				elements = append(elements, p.parseExpression())
			}
		}
		p.expectAndAdvance(lexer.RBRACKET)
		return &ast.List{Elements: elements}

	case lexer.LBRACE:
		// Dict literal
		p.advance()
		var keys []ast.Expr
		var values []ast.Expr
		if p.peek().Type != lexer.RBRACE {
			key := p.parseExpression()
			p.expectAndAdvance(lexer.COLON)
			value := p.parseExpression()
			keys = append(keys, key)
			values = append(values, value)
			for p.peek().Type == lexer.COMMA {
				p.advance()
				if p.peek().Type == lexer.RBRACE {
					break
				}
				key = p.parseExpression()
				p.expectAndAdvance(lexer.COLON)
				value = p.parseExpression()
				keys = append(keys, key)
				values = append(values, value)
			}
		}
		p.expectAndAdvance(lexer.RBRACE)
		return &ast.Dict{Keys: keys, Values: values}

	default:
		p.advance() // skip unknown token
		return nil
	}
}
