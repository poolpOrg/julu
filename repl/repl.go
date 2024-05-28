package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/poolpOrg/julu/lexer"
	"github.com/poolpOrg/julu/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, "julu>> ")
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

		io.WriteString(out, p.Parse().String()+"\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}
