package parser

import (
	"fmt"
	"strconv"

	"github.com/poolpOrg/julu/ast"
	"github.com/poolpOrg/julu/lexer"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.registerPrefix(lexer.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(lexer.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.FSTRING, p.parseFStringLiteral)
	p.registerPrefix(lexer.LOGICAL_NOT, p.parsePrefixExpression) // !x
	p.registerPrefix(lexer.BITWISE_NOT, p.parsePrefixExpression) // !x
	p.registerPrefix(lexer.SUB, p.parsePrefixExpression)         // -x
	p.registerPrefix(lexer.TRUE, p.parseBoolean)                 // true
	p.registerPrefix(lexer.FALSE, p.parseBoolean)                // false
	p.registerPrefix(lexer.IF, p.parseIfExpression)              // if (...)
	p.registerPrefix(lexer.FN, p.parseFunctionLiteral)           // fn (...
	p.registerPrefix(lexer.LEFT_PARENTHESIS, p.parseGroupedExpression)

	p.registerInfix(lexer.ADD, p.parseInfixExpression)
	p.registerInfix(lexer.SUB, p.parseInfixExpression)
	p.registerInfix(lexer.MUL, p.parseInfixExpression)
	p.registerInfix(lexer.DIV, p.parseInfixExpression)
	p.registerInfix(lexer.MOD, p.parseInfixExpression)

	p.registerInfix(lexer.BITWISE_AND, p.parseInfixExpression)
	p.registerInfix(lexer.BITWISE_OR, p.parseInfixExpression)
	p.registerInfix(lexer.BITWISE_XOR, p.parseInfixExpression)

	p.registerInfix(lexer.RSHIFT, p.parseInfixExpression)
	p.registerInfix(lexer.LSHIFT, p.parseInfixExpression)
	p.registerInfix(lexer.CIRCULAR_RSHIFT, p.parseInfixExpression)
	p.registerInfix(lexer.CIRCULAR_LSHIFT, p.parseInfixExpression)

	p.registerInfix(lexer.LOGICAL_AND, p.parseInfixExpression)
	p.registerInfix(lexer.LOGICAL_OR, p.parseInfixExpression)

	p.registerInfix(lexer.EQUALS, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQUALS, p.parseInfixExpression)
	p.registerInfix(lexer.LESSER_THAN, p.parseInfixExpression)
	p.registerInfix(lexer.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(lexer.LESSER_OR_EQUAL, p.parseInfixExpression)
	p.registerInfix(lexer.GREATER_OR_EQUAL, p.parseInfixExpression)
	p.registerInfix(lexer.LEFT_PARENTHESIS, p.parseCallExpression) // myFunction(x)

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.Lex()
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Parse() *ast.Program {
	program := ast.NewProgram()

	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := ast.NewLetStatement(p.curToken)

	if !p.expectPeek(lexer.IDENTIFIER) {
		return nil
	}
	stmt.Name = ast.NewIdentifier(p.curToken)

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	for !p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := ast.NewReturnStatement(p.curToken)

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := ast.NewExpressionStatement(p.curToken)

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return ast.NewIdentifier(p.curToken)
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return ast.NewStringLiteral(p.curToken)
}

func (p *Parser) parseFStringLiteral() ast.Expression {
	return ast.NewFStringLiteral(p.curToken)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return ast.NewIntegerLiteral(p.curToken, value)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := ast.NewPrefixExpression(p.curToken, p.curToken.Literal)

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.NewInfixExpression(p.curToken, left)

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return ast.NewBoolean(p.curToken, p.curTokenIs(lexer.TRUE))
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RIGHT_PARENTHESIS) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := ast.NewIfExpression(p.curToken)

	if !p.expectPeek(lexer.LEFT_PARENTHESIS) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RIGHT_PARENTHESIS) {
		return nil
	}

	if !p.expectPeek(lexer.LEFT_CURLY_BRACKET) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()
		if !p.expectPeek(lexer.LEFT_CURLY_BRACKET) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := ast.NewBlockStatement(p.curToken)

	p.nextToken()

	for !p.curTokenIs(lexer.RIGHT_CURLY_BRACKET) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	expression := ast.NewFunctionLiteral(p.curToken)

	if !p.expectPeek(lexer.LEFT_PARENTHESIS) {
		return nil
	}

	expression.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(lexer.LEFT_CURLY_BRACKET) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(lexer.RIGHT_PARENTHESIS) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := ast.NewIdentifier(p.curToken)
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := ast.NewIdentifier(p.curToken)
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(lexer.RIGHT_PARENTHESIS) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := ast.NewCallExpression(p.curToken, function)
	expression.Arguments = p.parseCallArguments()
	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(lexer.RIGHT_PARENTHESIS) {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RIGHT_PARENTHESIS) {
		return nil
	}

	return args
}
