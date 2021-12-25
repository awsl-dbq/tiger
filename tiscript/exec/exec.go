package exec

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/awsl-dbq/tiger/tiscript/evaluator"
	"github.com/awsl-dbq/tiger/tiscript/lexer"
	"github.com/awsl-dbq/tiger/tiscript/object"
	"github.com/awsl-dbq/tiger/tiscript/parser"
)

func Run(fileName string) {
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()
	f, e := os.Open(fileName)
	if e != nil {
		fmt.Printf("error %v \n", e)
		return
	}
	b, e := ioutil.ReadAll(f)
	if e != nil {
		fmt.Printf("error %v", e)
		return
	}

	l := lexer.New(string(b))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Printf("%v", msg)
		}
	}
	evaluator.DefineMacros(program, macroEnv)
	expanded := evaluator.ExpandMacros(program, macroEnv)
	evaluated := evaluator.Eval(expanded, env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}

}
