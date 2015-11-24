package parser

import (
	"errors"
	"github.com/suzuken/gs/lexer"
	"github.com/suzuken/gs/types"
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
	// start s-expression
	if p.lex.TokenText() == "(" {
		var list []types.Expression
		p.lex.Scan()
		// recursive scan until ')'
	LOOP:
		for {
			switch p.lex.TokenText() {
			case ")":
				break LOOP
			default:
				p.lex.Scan()
				list = append(list, p.lex.TokenText())
			}
		}
		return list, nil
	} else if p.lex.TokenText() == ")" {
		return exps, errors.New("unexpected ')'")
	} else {
		return p.lex.TokenText(), nil
	}
}
