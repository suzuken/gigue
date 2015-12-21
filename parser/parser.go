package parser

import (
	"errors"
	"github.com/suzuken/gigue/lexer"
	"github.com/suzuken/gigue/types"
	"strconv"
)

type Parser struct {
	lex *lexer.Lex
}

// New returns new parser
func New(lex *lexer.Lex) *Parser {
	return &Parser{lex}
}

// Parse scans given tokens and return it into expression.
func (p *Parser) Parse() (exps types.Expression, err error) {
	token, err := p.lex.Next()
	if err != nil {
		return nil, err
	}
	// start s-expression
	if token == "(" {
		var list []types.Expression
		// recursive scan until ")"
		for {
			if p.lex.Peek() == ')' {
				break
			}
			ex, err := p.Parse()
			if err != nil {
				return nil, err
			}
			list = append(list, ex)
		}
		// detect by Peek(), so scanner should read next rune.
		p.lex.Scan()
		return list, nil
	} else if token == ")" {
		return nil, errors.New("unexpected ')'")
	} else {
		if token == "#t" {
			return types.Boolean(true), nil
		} else if token == "#f" {
			return types.Boolean(false), nil
		}

		// if token is string, get unquoted.
		// It's get test from \"test\".
		if p.lex.IsTokenString() {
			str, err := strconv.Unquote(p.lex.TokenText())
			if err != nil {
				return "", err
			}
			return str, nil
		}

		// try conversion to float. if failed, deal with symbol.
		if n, err := strconv.ParseFloat(token, 64); err == nil {
			return types.Number(n), nil
		}

		return types.Symbol(token), nil
	}
}
