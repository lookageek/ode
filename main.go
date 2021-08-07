package main

import (
	"fmt"
	"lookageek.com/ode/repl"
	"os"
	"os/user"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", currentUser.Username)
	fmt.Printf("Feel free to type some commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
