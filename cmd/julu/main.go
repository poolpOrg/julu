package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/poolpOrg/julu/evaluator"
	"github.com/poolpOrg/julu/lexer"
	"github.com/poolpOrg/julu/object"
	"github.com/poolpOrg/julu/parser"
	"github.com/poolpOrg/julu/repl"
	"golang.org/x/term"
)

func main() {
	flag.Parse()

	if term.IsTerminal(int(os.Stdin.Fd())) {
		os.Exit(repl.Start(os.Stdin, os.Stdout))
	}

	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()
	for {
		scanned := scanner.Scan()
		if !scanned {
			os.Exit(0)
		}

		line := scanner.Text()
		l := lexer.New(bufio.NewReader(strings.NewReader(line)))
		p := parser.New(l)

		if len(p.Errors()) > 0 {
			printParserErrors(os.Stderr, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(p.Parse(), env)
		if evaluated != nil {
			io.WriteString(os.Stdout, evaluated.Inspect()+"\n")
		}
	}

}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
