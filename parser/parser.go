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

	entryPoint ast.Expression
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
	p.registerPrefix(lexer.NULL, p.parseNull)
	p.registerPrefix(lexer.FSTRING, p.parseFStringLiteral)
	p.registerPrefix(lexer.LOGICAL_NOT, p.parsePrefixExpression) // !x
	p.registerPrefix(lexer.BITWISE_NOT, p.parsePrefixExpression) // !x
	p.registerPrefix(lexer.SUB, p.parsePrefixExpression)         // -x
	p.registerPrefix(lexer.TRUE, p.parseBoolean)                 // true
	p.registerPrefix(lexer.FALSE, p.parseBoolean)                // false
	p.registerPrefix(lexer.IF, p.parseIfExpression)              // if (...)
	p.registerPrefix(lexer.FN, p.parseFunctionLiteral)           // fn (...
	p.registerPrefix(lexer.LEFT_PARENTHESIS, p.parseGroupedExpression)
	p.registerPrefix(lexer.LEFT_SQUARE_BRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.LEFT_CURLY_BRACKET, p.parseHashLiteral)
	p.registerPrefix(lexer.LOOP, p.parseLoopStatement)
	p.registerPrefix(lexer.WHILE, p.parseWhileStatement)
	p.registerPrefix(lexer.UNTIL, p.parseUntilStatement)
	p.registerPrefix(lexer.FOR, p.parseForStatement)
	p.registerPrefix(lexer.MATCH, p.parseMatchExpression)

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
	p.registerInfix(lexer.LEFT_SQUARE_BRACKET, p.parseIndexExpression)

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) pushError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	p.pushError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
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
	var ret ast.Statement
	switch p.curToken.Type {
	case lexer.LET:
		ret = p.parseLetStatement()
	case lexer.DONE:
		ret = p.parseDoneStatement()
	case lexer.RETURN:
		ret = p.parseReturnStatement()
	case lexer.BREAK:
		ret = p.parseBreakStatement()
	case lexer.CONTINUE:
		ret = p.parseContinueStatement()
	default:
		ret = p.parseExpressionStatement()
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return ret
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

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := ast.NewReturnStatement(p.curToken)

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	return ast.NewBreakStatement(p.curToken)
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	return ast.NewContinueStatement(p.curToken)
}

func (p *Parser) parseDoneStatement() *ast.DoneStatement {
	return ast.NewDoneStatement(p.curToken)
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
	p.pushError(fmt.Sprintf("no prefix parse function for %s found", t))
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
		p.pushError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
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

func (p *Parser) parseNull() ast.Expression {
	return ast.NewNull(p.curToken)
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

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) && !p.peekTokenIs(lexer.ARROW) {
		return nil
	}
	p.nextToken()
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()

		if p.peekTokenIs(lexer.IF) {
			p.nextToken()
			expression.ConditionalAlternative = p.parseIfExpression().(*ast.IfExpression)
			return expression
		}
		// if the next token is IS, then we have a conditional else

		if !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) && !p.peekTokenIs(lexer.ARROW) {
			return nil
		}
		p.nextToken()
		expression.Alternative = p.parseBlockStatement()
	}
	p.nextToken()

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := ast.NewBlockStatement(p.curToken)
	p.nextToken()

	if block.Token.Type == lexer.ARROW {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		return block
	}

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

	if p.peekTokenIs(lexer.IDENTIFIER) {
		p.nextToken()
		expression.Name = ast.NewIdentifier(p.curToken)
	}

	if p.peekTokenIs(lexer.LEFT_PARENTHESIS) {
		p.nextToken()
		expression.Parameters = p.parseFunctionParameters()
	}

	if !p.peekTokenIs(lexer.ARROW) && !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) {
		return nil
	}
	p.nextToken()

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
	expression.Parameters = p.parseCallParameters()
	return expression
}

func (p *Parser) parseCallParameters() []ast.Expression {
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

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := ast.NewArrayLiteral(p.curToken)
	array.Elements = p.parseExpressionList(lexer.RIGHT_SQUARE_BRACKET)
	return array
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := ast.NewIndexExpression(p.curToken, left)

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RIGHT_SQUARE_BRACKET) {
		return nil
	}

	return expression
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := ast.NewHashLiteral(p.curToken)
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(lexer.RIGHT_CURLY_BRACKET) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(lexer.RIGHT_CURLY_BRACKET) && !p.expectPeek(lexer.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(lexer.RIGHT_CURLY_BRACKET) {
		return nil
	}

	return hash
}

func (p *Parser) parseLoopStatement() ast.Expression {
	stmt := ast.NewLoopStatement(p.curToken)

	if !p.peekTokenIs(lexer.ARROW) && !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) {
		return nil
	}
	p.nextToken()

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() ast.Expression {
	stmt := ast.NewLoopStatement(p.curToken)

	p.nextToken()
	stmt.WhileCondition = p.parseExpression(LOWEST)

	if !p.peekTokenIs(lexer.ARROW) && !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) {
		return nil
	}
	p.nextToken()

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseUntilStatement() ast.Expression {
	stmt := ast.NewLoopStatement(p.curToken)

	p.nextToken()
	stmt.UntilCondition = p.parseExpression(LOWEST)

	if !p.peekTokenIs(lexer.ARROW) && !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) {
		return nil
	}
	p.nextToken()

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() ast.Expression {
	stmt := ast.NewForStatement(p.curToken)
	p.nextToken()

	stmt.Variable = p.parseIdentifier()

	if !p.peekTokenIs(lexer.IN) {
		return nil
	}
	p.nextToken()
	p.nextToken()

	stmt.Iterable = p.parseExpression(LOWEST)

	if !p.peekTokenIs(lexer.ARROW) && !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) {
		return nil
	}
	p.nextToken()

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseMatchExpression() ast.Expression {
	expression := ast.NewMatchExpression(p.curToken)

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) && !p.peekTokenIs(lexer.ARROW) {
		return nil
	}
	p.nextToken()

	expression.MatchBlock = p.parseMatchBlockStatement()

	for !p.curTokenIs(lexer.RIGHT_CURLY_BRACKET) {
		p.nextToken()
	}

	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()
		if !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) && !p.peekTokenIs(lexer.ARROW) {
			return nil
		}
		p.nextToken()
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseMatchBlockStatement() *ast.MatchBlockStatement {
	block := ast.NewMatchBlockStatement(p.curToken)
	p.nextToken()

	if block.Token.Type == lexer.ARROW {
		stmt := p.parseCaseExpression()
		if stmt == nil {
			p.pushError("expected if expression")
			return nil
		}

		if stmt, ok := stmt.(*ast.CaseExpression); !ok {
			p.pushError("expected if expression")
			return nil
		} else {
			block.Cases = append(block.Cases, *stmt)
		}

		return block
	}

	for !p.curTokenIs(lexer.RIGHT_CURLY_BRACKET) {
		stmt := p.parseCaseExpression()
		if stmt == nil {
			p.pushError("expected if expression")
			return nil
		}

		if stmt, ok := stmt.(*ast.CaseExpression); !ok {
			p.pushError("expected if expression")
			return nil
		} else {
			block.Cases = append(block.Cases, *stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseCaseExpression() ast.Expression {
	if !p.curTokenIs(lexer.CASE) {
		return nil
	}
	p.nextToken()

	expression := ast.NewCaseExpression(p.curToken)
	expression.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.IF) {
		p.nextToken()
		p.nextToken()
		expression.Guard = p.parseExpression(LOWEST)
	}

	if !p.peekTokenIs(lexer.LEFT_CURLY_BRACKET) && !p.peekTokenIs(lexer.ARROW) {
		return nil
	}
	p.nextToken()
	expression.Consequence = p.parseBlockStatement()
	return expression
}
