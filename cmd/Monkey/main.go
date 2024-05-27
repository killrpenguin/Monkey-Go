package main

import (
	"Monkey/internal/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! Welcome to Monkey-Lang!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
