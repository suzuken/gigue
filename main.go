// GS is exprerimental implementation of R5RS Scheme.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/suzuken/gs/eval"
	"github.com/suzuken/gs/lexer"
	"github.com/suzuken/gs/parser"
	"github.com/suzuken/gs/types"
	"io"
	"os"
	"strings"
)

func main() {
	env := eval.NewEnv()
	env.Setup()

	flag.Parse()
	if len(flag.Args()) == 0 {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "input error: %s\n", err)
			}
			if input == "exit\n" {
				return
			}
			exp, err := ev(strings.NewReader(input), env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "eval error: %s\n", err)
			} else {
				showEval(exp)
			}
		}
	} else {
		filename := flag.Arg(1)
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open file: %s\n", err)
			os.Exit(1)
		}
		defer f.Close()
		exp, err := ev(f, env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "eval error: %s\n", err)
			os.Exit(1)
		}
		showEval(exp)
	}
}

func showEval(exp types.Expression) {
	fmt.Printf("%v\n", exp)
}

func ev(r io.Reader, env *eval.Env) (types.Expression, error) {
	l := lexer.New()
	l.Init(r)
	l.Scan()
	p := parser.New(l)
	exps, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return eval.Eval(exps, env)
}
