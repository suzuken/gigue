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
		// recursive scan until ')'
		for {
			p.lex.Scan()
			if p.lex.TokenText() == ")" {
				break
			}
			ex, err := p.Parse()
			if err != nil {
				return nil, err
			}
			list = append(list, ex)
		}
		return list, nil
	} else if p.lex.TokenText() == ")" {
		return exps, errors.New("unexpected ')'")
	} else {
		return p.lex.TokenText(), nil
	}
}
