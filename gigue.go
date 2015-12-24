package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/suzuken/gigue/eval"
	"github.com/suzuken/gigue/types"
	"io"
	"os"
	"strings"
)

func main() {
	env := eval.NewEnv()
	env.Setup()

	flag.Parse()
	if flag.NArg() == 0 {
		var buf bytes.Buffer
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if buf.Len() == 0 {
				fmt.Print("> ")
			} else {
				fmt.Print(">> ")
			}
			if !scanner.Scan() {
				if scanner.Err() != nil {
					fmt.Fprintf(os.Stderr, "input error: %s\n", scanner.Err())
				}
				break
			}
			input := scanner.Text()
			if input == "exit" {
				return
			}
			s := strings.TrimSpace(input)
			if buf.Len() > 0 && s != ")" {
				buf.WriteRune(' ')
			} else {
				input = s
			}
			buf.WriteString(input)
			exp, err := eval.EvalReader(strings.NewReader(buf.String()), env)
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(os.Stderr, "eval error: %s\n", err)
					buf.Reset()
				}
			} else {
				output(exp)
				buf.Reset()
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
