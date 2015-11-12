package lexer

import (
	"errors"
	"fmt"
	"io"
	"math/big"
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

type Tokens []string

// Scan starts scan the whole program and return tokens
func (lex *Lex) Scan() (tokens Tokens, err error) {
	// start s-expression
	if lex.Token == '(' {
		lex.NextToken()
		// recursive scan
		ts, err := lex.Scan()
		if err != nil {
			return tokens, err
		}
		for _, t := range ts {
			tokens = append(tokens, t)
		}
		if lex.Token != ')' {
			return tokens, fmt.Errorf("')' is expected but get other one. failed. token: %s", lex.TokenText())
		}
		return tokens, nil
	} else if lex.Token == ')' {
		return tokens, errors.New("unexpected ')'")
	}
	return tokens, nil
}
