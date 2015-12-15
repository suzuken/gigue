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
	"text/scanner"
)

func main() {
	env := eval.NewEnv()
	env.Setup()

	flag.Parse()
	if flag.NArg() == 0 {
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
				output(exp)
			}
		}
	} else {
		if _, err := evalFile(flag.Arg(0), env); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
}

// output simply output results.
func output(exp types.Expression) {
	if exp != nil {
		fmt.Printf("%v\n", exp)
	}
}

func evalFile(filename string, env *eval.Env) (types.Expression, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ev(f, env)
}

// ev evaluate scheme program from io.Reader
func ev(r io.Reader, env *eval.Env) (types.Expression, error) {
	l := lexer.New()
	l.Init(r)
	l.Scan()
	p := parser.New(l)
	for l.Token != scanner.EOF {
		// fmt.Println(l.TokenText())
		exps, err := p.Parse()
		if err != nil {
			return nil, err
		}
		// fmt.Printf("%#v\n", exps)
		if _, err := eval.Eval(exps, env); err != nil {
			return nil, err
		}
		l.Scan()
	}
	return nil, nil
}
