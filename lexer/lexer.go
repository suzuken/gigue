package lexer

import (
	"errors"
	"fmt"
	"github.com/suzuken/gs/types"
	"io"
	"text/scanner"
)

type Lex struct {
	*scanner.Scanner
	Token rune
}

// New returns new lexer
func New() *Lex {
	var s scanner.Scanner
	// only scan characters. implement lexer myself.
	s.Mode &^= scanner.ScanChars | scanner.ScanRawStrings
	return &Lex{
		Scanner: &s,
	}
}

// Init initialize lexer
func (l *Lex) Init(r io.Reader) {
	l.Scanner.Init(r)
}

// NextToken gets next token to lexer
func (lex *Lex) NextToken() {
	lex.Token = lex.Scanner.Scan()
}

// Error creates error including current token context.
func (lex *Lex) Error(msg string) error {
	return fmt.Errorf("%s: %v", msg, lex.Token)
}

// Scan starts scan the whole program and return tokens
func (lex *Lex) Scan() (exps types.Expression, err error) {
	// start s-expression
	if lex.Token == '(' {
		var list []types.Expression
		lex.NextToken()
		// recursive scan until ')'
	LOOP:
		for {
			switch lex.Token {
			case ')':
				break LOOP
			default:
				ts, err := lex.Scan()
				if err != nil {
					return exps, err
				}
				lex.NextToken()
				list = append(list, ts)
			}
		}
		return list, nil
	} else if lex.Token == ')' {
		return exps, errors.New("unexpected ')'")
	} else {
		return lex.TokenText(), nil
	}
}
