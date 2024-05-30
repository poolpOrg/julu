package lexer

import (
	"bufio"
	"io"
	"unicode"
)

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	ADD  = "ADD"
	SUB  = "SUB"
	MUL  = "MUL"
	DIV  = "DIV"
	MOD  = "MOD"
	INCR = "INCR"
	DECR = "DECR"

	LESSER_THAN      = "LESSER_THAN"
	GREATER_THAN     = "GREATER_THAN"
	LESSER_OR_EQUAL  = "LESSER_OR_EQUAL"
	GREATER_OR_EQUAL = "GREATER_OR_EQUAL"
	EQUALS           = "EQUALS"
	NOT_EQUALS       = "NOT_EQUALS"

	BITWISE_AND     = "BITWISE_AND"
	BITWISE_OR      = "BITWISE_OR"
	BITWISE_XOR     = "BITWISE_XOR"
	BITWISE_NOT     = "BITWISE_NOT"
	LSHIFT          = "LSHIFT"
	RSHIFT          = "RSHIFT"
	CIRCULAR_LSHIFT = "CIRCULAR_LSHIFT"
	CIRCULAR_RSHIFT = "CIRCULAR_RSHIFT"

	ASSIGN         = "ASSIGN"
	ADD_AND_ASSIGN = "ADD_AND_ASSIGN"
	SUB_AND_ASSIGN = "SUB_AND_ASSIGN"
	MUL_AND_ASSIGN = "MUL_AND_ASSIGN"
	DIV_AND_ASSIGN = "DIV_AND_ASSIGN"
	MOD_AND_ASSIGN = "MOD_AND_ASSIGN"

	BITWISE_AND_ASSIGN = "BITWISE_AND_ASSIGN"
	BITWISE_OR_ASSIGN  = "BITWISE_OR_ASSIGN"
	BITWISE_XOR_ASSIGN = "BITWISE_XOR_ASSIGN"
	LSHIFT_ASSIGN      = "LSHIFT_ASSIGN"
	RSHIFT_ASSIGN      = "RSHIFT_ASSIGN"

	LOGICAL_AND = "LOGICAL_AND"
	LOGICAL_OR  = "LOGICAL_OR"
	LOGICAL_NOT = "LOGICAL_NOT"

	RUNE       = "RUNE"
	STRING     = "STRING"
	FSTRING    = "FSTRING"
	INTEGER    = "INTEGER"
	FLOAT      = "FLOAT"
	IDENTIFIER = "IDENTIFIER"

	ARROW                = "ARROW"
	LEFT_PARENTHESIS     = "LEFT_PARENTHESIS"
	RIGHT_PARENTHESIS    = "RIGHT_PARENTHESIS"
	LEFT_CURLY_BRACKET   = "LEFT_CURLY_BRACKET"
	RIGHT_CURLY_BRACKET  = "RIGHT_CURLY_BRACKET"
	LEFT_SQUARE_BRACKET  = "LEFT_SQUARE_BRACKET"
	RIGHT_SQUARE_BRACKET = "RIGHT_SQUARE_BRACKET"

	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"
	COLON     = "COLON"
	DOT       = "DOT"

	IS = "IS"
	IN = "IN"

	// Keywords
	NULL  = "NULL"
	TRUE  = "TRUE"
	FALSE = "FALSE"

	// Control flow
	IF     = "IF"
	ELSE   = "ELSE"
	MATCH  = "MATCH"
	CASE   = "CASE"
	RETURN = "RETURN"

	LET = "LET"

	FN = "FN"

	FOR      = "FOR"
	LOOP     = "LOOP"
	WHILE    = "WHILE"
	UNTIL    = "UNTIL"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	DONE     = "DONE"
)

var keywords = map[string]TokenType{
	"let": LET,

	"null": NULL,

	"is":  IS,
	"in":  IN,
	"and": LOGICAL_AND,
	"or":  LOGICAL_OR,
	"not": LOGICAL_NOT,

	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"fn":     FN,

	"loop":     LOOP,
	"while":    WHILE,
	"until":    UNTIL,
	"for":      FOR,
	"match":    MATCH,
	"case":     CASE,
	"break":    BREAK,
	"continue": CONTINUE,
	"done":     DONE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

type Position struct {
	line   int
	column int
}

func (p *Position) Line() int {
	return p.line
}

func (p *Position) Column() int {
	return p.column
}

type Token struct {
	Type     TokenType
	Literal  string
	position *Position
}

func (t *Token) String() string {
	return t.Literal
}

func (t *Token) Position() *Position {
	return t.position
}

func tokenFromLexer(t TokenType, position Position, literal string) Token {
	return Token{
		Type:     t,
		Literal:  literal,
		position: &position,
	}
}

type Lexer struct {
	reader *bufio.Reader
	pos    Position
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	l.pos.column--
}

func (l *Lexer) skipLine() {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if r == '\n' {
			l.backup()
			break
		}
	}
}

func (l *Lexer) skipMultilineComment() {
	var lastRune rune
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if r == '\n' {
			l.pos.column = 0
			l.pos.line++
		} else if lastRune == '*' && r == '/' {
			break
		} else {
			l.pos.column++
		}
		lastRune = r
	}
}

func New(reader *bufio.Reader) *Lexer {
	return &Lexer{
		reader: reader,
		pos:    Position{line: 1, column: 0},
	}
}

func (l *Lexer) Lex() Token {
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return tokenFromLexer(EOF, l.pos, "")
			}
			panic(err)
		}
		l.pos.column++

		//will be used as start position for multi-rune tokens
		startPos := l.pos
		switch r {
		case '\n':
			l.pos.column = 0
			l.pos.line++

		case '#':
			l.skipLine()

		case ',':
			return tokenFromLexer(COMMA, startPos, string(r))
		case ';':
			return tokenFromLexer(SEMICOLON, startPos, string(r))
		case ':':
			return tokenFromLexer(COLON, startPos, string(r))
		case '.':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if unicode.IsDigit(nextR) {
					l.backup()
					tokenType, lit := l.lexNumber(FLOAT)
					return tokenFromLexer(tokenType, startPos, "."+lit)
				}
				l.backup()
			}
			return tokenFromLexer(DOT, startPos, string(r))

		case '(':
			return tokenFromLexer(LEFT_PARENTHESIS, startPos, string(r))
		case ')':
			return tokenFromLexer(RIGHT_PARENTHESIS, startPos, string(r))
		case '{':
			return tokenFromLexer(LEFT_CURLY_BRACKET, startPos, string(r))
		case '}':
			return tokenFromLexer(RIGHT_CURLY_BRACKET, startPos, string(r))
		case '[':
			return tokenFromLexer(LEFT_SQUARE_BRACKET, startPos, string(r))
		case ']':
			return tokenFromLexer(RIGHT_SQUARE_BRACKET, startPos, string(r))

		case '+':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '+' {
					return tokenFromLexer(INCR, startPos, "++")
				} else if nextR == '=' {
					return tokenFromLexer(ADD_AND_ASSIGN, startPos, "+=")
				}
				l.backup()
			}
			return tokenFromLexer(ADD, startPos, string(r))

		case '-':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '-' {
					return tokenFromLexer(DECR, startPos, "--")
				} else if nextR == '=' {
					return tokenFromLexer(SUB_AND_ASSIGN, startPos, "-=")
				}
				l.backup()
			}
			return tokenFromLexer(SUB, startPos, string(r))

		case '*':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(MUL_AND_ASSIGN, startPos, "*=")
				}
				l.backup()
			}
			return tokenFromLexer(MUL, startPos, string(r))

		case '/':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(DIV_AND_ASSIGN, startPos, "/=")
				} else if nextR == '/' {
					l.skipLine()
					continue
				} else if nextR == '*' {
					l.skipMultilineComment()
					continue
				}
				l.backup()
			}
			return tokenFromLexer(DIV, startPos, string(r))

		case '%':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(MOD_AND_ASSIGN, startPos, "%=")
				}
				l.backup()
			}
			return tokenFromLexer(MOD, startPos, string(r))

		case '=':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(EQUALS, startPos, "==")
				} else if nextR == '>' {
					return tokenFromLexer(ARROW, startPos, "=>")
				}
				l.backup()
			}
			return tokenFromLexer(ASSIGN, startPos, string(r))

		case '&':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(BITWISE_AND_ASSIGN, startPos, "&=")
				} else if nextR == '&' {
					return tokenFromLexer(LOGICAL_AND, startPos, "&&")
				}
				l.backup()
			}
			return tokenFromLexer(BITWISE_AND, startPos, string(r))

		case '|':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(BITWISE_OR_ASSIGN, startPos, "|=")
				} else if nextR == '|' {
					return tokenFromLexer(LOGICAL_OR, startPos, "||")
				}
				l.backup()
			}
			return tokenFromLexer(BITWISE_OR, startPos, string(r))

		case '^':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(BITWISE_XOR_ASSIGN, startPos, "^=")
				}
				l.backup()
			}
			return tokenFromLexer(BITWISE_XOR, startPos, string(r))

		case '~':
			return tokenFromLexer(BITWISE_NOT, startPos, string(r))

		case '<':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(LESSER_OR_EQUAL, startPos, "<=")
				} else if nextR == '<' {
					nextR, _, err := l.reader.ReadRune()
					if err == nil {
						l.pos.column++
						if nextR == '=' {
							return tokenFromLexer(LSHIFT_ASSIGN, startPos, "<<=")
						} else if nextR == '<' {
							return tokenFromLexer(CIRCULAR_LSHIFT, startPos, "<<<")
						}
						l.backup()
					}
					return tokenFromLexer(LSHIFT, startPos, "<<")
				}
				l.backup()
			}
			return tokenFromLexer(LESSER_THAN, startPos, string(r))

		case '>':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(GREATER_OR_EQUAL, startPos, ">=")
				} else if nextR == '>' {
					nextR, _, err := l.reader.ReadRune()
					if err == nil {
						l.pos.column++
						if nextR == '=' {
							return tokenFromLexer(RSHIFT_ASSIGN, startPos, ">>=")
						} else if nextR == '>' {
							return tokenFromLexer(CIRCULAR_RSHIFT, startPos, ">>>")
						}
						l.backup()
					}
					return tokenFromLexer(RSHIFT, startPos, ">>")
				}
				l.backup()
			}
			return tokenFromLexer(GREATER_THAN, startPos, string(r))

		case '!':
			nextR, _, err := l.reader.ReadRune()
			if err == nil {
				l.pos.column++
				if nextR == '=' {
					return tokenFromLexer(NOT_EQUALS, startPos, "!=")
				}
				l.backup()
			}
			return tokenFromLexer(LOGICAL_NOT, startPos, string(r))

		case '\'':
			l.backup()
			tokenType, lit := l.lexRune()
			if tokenType != RUNE {
				return tokenFromLexer(ILLEGAL, startPos, lit)
			}
			return tokenFromLexer(tokenType, startPos, lit)

		case '"':
			l.backup()
			tokenType, lit := l.lexString(false, '"')
			if tokenType != STRING {
				return tokenFromLexer(ILLEGAL, startPos, lit)
			}
			return tokenFromLexer(tokenType, startPos, lit)

		case '`':
			l.backup()
			tokenType, lit := l.lexString(true, '`')
			if tokenType != STRING {
				return tokenFromLexer(ILLEGAL, startPos, lit)
			}
			return tokenFromLexer(tokenType, startPos, lit)

		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) {
				l.backup()
				tokenType, lit := l.lexNumber(INTEGER)
				return tokenFromLexer(tokenType, startPos, lit)
			} else if unicode.IsLetter(r) || r == '_' {
				l.backup()
				tokenType, lit := l.lexIdentifier()
				if tokenType == FSTRING {
					return tokenFromLexer(tokenType, startPos, lit)
				}
				if tokenType != IDENTIFIER {
					return tokenFromLexer(ILLEGAL, startPos, lit)
				}
				return tokenFromLexer(LookupIdent(lit), startPos, lit)
			}

			return tokenFromLexer(ILLEGAL, startPos, string(r))
		}
	}
}

// XXX - fix to handle floats and other bases than 10
func (l *Lexer) lexNumber(tokenType TokenType) (TokenType, string) {
	var lit string
	base := 10
	baseOffset := 0
	offset := 0

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return tokenType, lit
			}
		}

		l.pos.column++
		if r == '.' && tokenType == INTEGER {
			tokenType = FLOAT
			base = 10
			baseOffset = offset + 1
			lit = lit + string(r)
		} else if unicode.IsDigit(r) {
			if base == 2 && !(r >= '0' && r <= '1') {
				l.backup()
				return tokenType, lit
			}
			if base == 8 && !(r >= '0' && r <= '7') {
				l.backup()
				return tokenType, lit
			}
			lit = lit + string(r)
		} else if r == '_' {
		} else {
			if len(lit)-baseOffset == 1 && lit[len(lit)-baseOffset-1] == '0' && (r == 'b' || r == 'd' || r == 'o' || r == 'x') {
				if r == 'b' {
					base = 2
				} else if r == 'o' {
					base = 8
				} else if r == 'x' {
					base = 16
				}
				lit = lit + string(r)
			} else {
				if base == 16 && ((r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
					lit = lit + string(r)
				} else {
					// scanned something not in the integer
					l.backup()
					return tokenType, lit
				}
			}
		}
		offset += 1
	}
}

func (l *Lexer) lexIdentifier() (TokenType, string) {
	var lit string
	var idx int
	var fbyte rune

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		l.pos.column++

		if idx == 0 {
			fbyte = r
		}

		if r == '"' {
			l.backup()
			if fbyte == 'f' && idx == 1 {
				tokenType, lit := l.lexString(false, '"')
				if tokenType != STRING {
					return ILLEGAL, lit
				}
				return FSTRING, lit
			}
		}

		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			l.backup()
			break
		}
		lit += string(r)
		idx++
	}

	return IDENTIFIER, lit
}

func (l *Lexer) lexRune() (TokenType, string) {
	var lit string

	r, _, err := l.reader.ReadRune()
	if err != nil {
		panic(err)
	}
	if r != '\'' {
		panic("expected a double quote")
	}
	l.pos.column++

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return EOF, ""
			}
			panic(err)
		}
		l.pos.column++

		if r == '\'' {
			return RUNE, lit
		}

		if r == '\\' {
			// XXX - handle escape sequences
		}

		if len(lit) > 1 {
			return ILLEGAL, lit
		}
		lit += string(r)
	}
}

func (l *Lexer) lexString(raw bool, delimiter rune) (TokenType, string) {
	var lit string

	r, _, err := l.reader.ReadRune()
	if err != nil {
		panic(err)
	}
	if r != delimiter {
		panic("something went wrong, expected a double-quote or back-tick")
	}
	l.pos.column++

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return EOF, ""
			}
			panic(err)
		}
		l.pos.column++

		if r == delimiter {
			return STRING, lit
		}

		if r == '\\' && !raw {
			r, _, err := l.reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return EOF, ""
				}
				panic(err)
			}
			l.pos.column++

			switch r {
			case 'n':
				lit += "\n"
			case 't':
				lit += "\t"
			case 'r':
				lit += "\r"
			case 'b':
				lit += "\b"
			case 'f':
				lit += "\f"
			case '\\':
				lit += "\\"
			case '"':
				lit += "\""
			case '`':
				lit += "`"
			default:
				lit += "\\" + string(r) // If it's not a recognized escape sequence, keep original
			}
		} else {
			lit += string(r)
		}
	}
}
