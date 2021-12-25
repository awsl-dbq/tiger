package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/awsl-dbq/tiger/tiscript/exec"
	"github.com/awsl-dbq/tiger/tiscript/repl"
)

var f = flag.String("f", "", "fileName")
var r = flag.Bool("r", false, "start repl")

func main() {
	flag.Parse()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	if *f != "" {
		exec.Run(*f)
	}
	if *r {
		fmt.Printf("Hello %s!\nWelcome to TiScript REPL!\n", user.Username)
		repl.Start(os.Stdin, os.Stdout)
	}
}
