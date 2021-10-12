package repl

import (
	"bufio"
	"fmt"
	"io"

	"lookageek.com/ode/evaluator"
	"lookageek.com/ode/lexer"
	"lookageek.com/ode/parser"
)

const PROMPT = ">> "

// REPL Start function is an endless loop waiting on input
// at the terminal and press of enter key
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		// wait for code to be entered in the terminal
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// when code is entered and enter is pressed, start the
		// processing of that value
		line := scanner.Text()
		lex := lexer.New(line)
		p := parser.New(lex)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
