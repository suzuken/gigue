package lexer

import (
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

// Scan makes scan
func (lex *Lex) Scan() {
	lex.Token = lex.Scanner.Scan()
}
