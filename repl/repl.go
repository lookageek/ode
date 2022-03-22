package repl

import (
	"bufio"
	"fmt"
	"io"

	"lookageek.com/ode/compiler"
	"lookageek.com/ode/evaluator"
	"lookageek.com/ode/lexer"
	"lookageek.com/ode/object"
	"lookageek.com/ode/parser"
	"lookageek.com/ode/vm"
)

const PROMPT = ">> "

// REPL Start function is an endless loop waiting on input
// at the terminal and press of enter key
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

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

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func StartVm(in io.Reader, out io.Writer) {
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

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n%s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n%s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
