package main

import (
	"fmt"
	"os"
	"os/user"

	"lookageek.com/ode/repl"
)

// main method is the entrypoint for the REPL interface
func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Ode programming language!\n", currentUser.Username)
	fmt.Printf("Feel free to type some commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
