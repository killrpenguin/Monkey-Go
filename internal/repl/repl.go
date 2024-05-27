package repl

import (
	"bufio"
	"fmt"
	"io"
	"Monkey/internal/lexer"
	"Monkey/internal/token"
)

const PROMPT = "Monkey DO >> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		lexer := lexer.New(line)

		for tkn := lexer.NextToken(); tkn.Type != token.EOF; tkn = lexer.NextToken() {
			fmt.Fprintf(out, "%+v\n", tkn)
		}
	}
}
