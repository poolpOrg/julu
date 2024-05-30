package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/poolpOrg/julu/evaluator"
	"github.com/poolpOrg/julu/lexer"
	"github.com/poolpOrg/julu/object"
	"github.com/poolpOrg/julu/parser"
	"github.com/poolpOrg/julu/repl"
	"golang.org/x/term"
)

func main() {
	var opt_mode string
	flag.StringVar(&opt_mode, "mode", "", "mode to run the interpreter in")
	flag.Parse()

	if term.IsTerminal(int(os.Stdin.Fd())) && flag.NArg() == 0 {
		os.Exit(repl.Start(os.Stdin, os.Stdout))
	}

	var err error
	var input io.Reader = os.Stdin
	if flag.NArg() != 0 {
		input, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open file: %s\n", err)
			os.Exit(1)
		}
	}

	if opt_mode == "lexer" {
		l := lexer.New(bufio.NewReader(input))
		for tok := l.Lex(); tok.Type != lexer.EOF; tok = l.Lex() {
			switch tok.Type {
			case lexer.STRING, lexer.FSTRING, lexer.IDENTIFIER:
				fmt.Printf("[%d:%d] %s => %s\n",
					tok.Position().Line(), tok.Position().Column(), tok.Type, tok.Literal)
			default:
				fmt.Printf("[%d:%d] %s\n",
					tok.Position().Line(), tok.Position().Column(), tok.Type)
			}
		}
		os.Exit(0)
	}

	if opt_mode == "ast" {
		l := lexer.New(bufio.NewReader(input))
		p := parser.New(l)
		program := p.Parse()
		if p.Errors() != nil {
			printParserErrors(os.Stderr, p.Errors())
		}
		fmt.Println(program.Inspect())
		os.Exit(0)
	}

	env := object.NewEnvironment()

	code, err := io.ReadAll(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read input: %s\n", err)
		os.Exit(1)
	}

	l := lexer.New(bufio.NewReader(bytes.NewReader(code)))
	p := parser.New(l)
	program := p.Parse()
	if program == nil {
		os.Exit(1)
	}

	if p.Errors() != nil {
		printParserErrors(os.Stderr, p.Errors())
	}

	evaluated := evaluator.Eval(program, env)
	if entryPoint, ok := env.Get("main"); ok {
		evaluated = evaluator.EvalFunctionObject(entryPoint, env)
	}

	if evaluated != nil {
		if evaluated.Type() == object.VOID_OBJ {
			os.Exit(0)
		}
		if evaluated.Type() == object.INTEGER_OBJ {
			os.Exit(int(evaluated.(*object.Integer).Value))
		}
		if evaluated.Type() == object.BOOLEAN_OBJ {
			if evaluated.(*object.Boolean).Value {
				os.Exit(0)
			}
			os.Exit(1)
		}
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Fprintf(os.Stderr, "error: %s\n", evaluated.Inspect())
			os.Exit(1)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
