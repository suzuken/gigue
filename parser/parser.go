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

// parseList is helper for parsing between ( and ).
// return tokens.
func (p *Parser) parseList() ([]types.Expression, error) {
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
}

// Parse scans given tokens and return it into expression.
func (p *Parser) Parse() (exps types.Expression, err error) {
	token, err := p.lex.Next()
	if err != nil {
		return nil, err
	}
	if token == "'" {
		// if start with (, deal as list.
		// recursive scan until ")"
		// if start with "'(", it's list.
		if p.lex.Peek() == '(' {
			// this is (. skip it.
			if _, err := p.lex.Next(); err != nil {
				return nil, err
			}
			tokens, err := p.parseList()
			if err != nil {
				return nil, err
			}
			return types.NewList(tokens...), nil
		}
		// if not start with (, it's simply string (return string itself).
		return p.lex.Next()
	}

	// start s-expression
	if token == "(" {
		return p.parseList()
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
			return strconv.Unquote(p.lex.TokenText())
		}

		// try conversion to float. if failed, deal with symbol.
		if n, err := strconv.ParseFloat(token, 64); err == nil {
			return types.Number(n), nil
		}

		return types.Symbol(token), nil
	}
}
