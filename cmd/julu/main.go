package main

import (
	"bufio"
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

	l := lexer.New(bufio.NewReader(input))
	env := object.NewEnvironment()
	p := parser.New(l)

	if len(p.Errors()) > 0 {
		printParserErrors(os.Stderr, p.Errors())
	}

	evaluated := evaluator.Eval(p.Parse(), env)
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.Inspect()+"\n")
	}

}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
