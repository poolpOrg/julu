package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/poolpOrg/julu/evaluator"
	"github.com/poolpOrg/julu/lexer"
	"github.com/poolpOrg/julu/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, ">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(bufio.NewReader(strings.NewReader(line)))
		p := parser.New(l)

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(p.Parse())
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect()+"\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
