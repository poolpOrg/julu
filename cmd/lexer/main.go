package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/poolpOrg/julu/lexer"
)

func main() {
	l := lexer.New(bufio.NewReader(os.Stdin))
	for {
		t := l.NextToken()
		if t.Type == lexer.EOF {
			break
		}
		fmt.Printf("[%d:%d] %s (\"%s\")\n",
			t.Position().Line(), t.Position().Column(), t.Type, t.Literal)
	}
}
