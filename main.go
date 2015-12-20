// GS is exprerimental implementation of R5RS Scheme.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/suzuken/gs/eval"
	"github.com/suzuken/gs/types"
	"os"
	"strings"
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
			exp, err := eval.EvalReader(strings.NewReader(input), env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "eval error: %s\n", err)
			} else {
				output(exp)
			}
		}
	} else {
		if _, err := eval.EvalFile(flag.Arg(0), env); err != nil {
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
