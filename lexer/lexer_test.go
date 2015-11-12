package lexer

import (
	"github.com/suzuken/gs/types"
	"reflect"
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	lex := New()
	r := strings.NewReader("(print 1)")
	actual := []types.Expression{
		"print",
		"1",
	}
	lex.Init(r)
	lex.NextToken()
	tokens, err := lex.Scan()
	if err != nil {
		t.Fatalf("lexer failed: %s", err)
	}
	if !reflect.DeepEqual(tokens, actual) {
		t.Fatalf("tokens is not expected. %v", tokens)
	}
}

func TestLexRecursive(t *testing.T) {
	lex := New()
	r := strings.NewReader("(define (square x) (* x x))")
	expected := []types.Expression{
		"define",
		[]types.Expression{
			"square",
			"x",
		},
		[]types.Expression{
			"*",
			"x",
			"x",
		},
	}
	lex.Init(r)
	lex.NextToken()
	tokens, err := lex.Scan()
	if err != nil {
		t.Fatalf("lexer failed: %s", err)
	}
	if !reflect.DeepEqual(tokens, expected) {
		t.Fatalf("tokens is not expected. %v", tokens)
	}
}
