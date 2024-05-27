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
		fmt.Println(t.Type, t.Position(), t.Literal)
	}
}
