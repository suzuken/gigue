package lexer

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"
)

// Lex is lexer for gigue.
// Lex provides token for scheme.
type Lex struct {
	*scanner.Scanner
	Token rune
}

// New returns new lexer
func New() *Lex {
	var s scanner.Scanner
	// only scan characters. implement lexer myself.
	// s.Mode &^= scanner.ScanChars | scanner.ScanRawStrings
	s.Mode &^= scanner.ScanChars
	// s.Mode &^= scanner.ScanRawStrings
	return &Lex{
		Scanner: &s,
	}
}

// Next returns next scheme token
// Parser calls this Next() iteratively for fetching token.
func (l *Lex) Next() (string, error) {
	l.Scan()
	switch l.Token {
	case '#':
		// convert boolean. text/scanner split # as 1 token.
		switch peek := l.Peek(); peek {
		case 't', 'f':
			l.Scan()
			return fmt.Sprintf("#%s", l.TokenText()), nil
		default:
			return "", errors.New("unknown hash symbol")
		}
	case '(', ')':
		return l.TokenText(), nil
	default:
		// concatenate tokens for symbol until end of characters.
		token := l.TokenText()
		for {
			r := l.Peek()
			if unicode.IsSpace(r) || strings.ContainsRune("()'.,;", r) || r == scanner.EOF {
				break
			}
			token = fmt.Sprintf("%s%c", token, r)
			l.Scanner.Next()
		}
		return token, nil
	}
}

// IsTokenString check if token is string.
// It's use for get unquoted string in parser.
func (l *Lex) IsTokenString() bool {
	return l.Token == scanner.String
}

// TokenAll returns all token generated by lexer.
// It is useful for testing lexer.
func (l *Lex) TokenAll() (tokens []string, err error) {
	for l.Peek() != scanner.EOF {
		token, err := l.Next()
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, token)
	}
	return
}

// Init initialize lexer
func (l *Lex) Init(r io.Reader) {
	l.Scanner.Init(r)
}

// Scan makes scan
func (lex *Lex) Scan() {
	lex.Token = lex.Scanner.Scan()
}
