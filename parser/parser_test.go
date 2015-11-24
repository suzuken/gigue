package parser

import (
	"github.com/suzuken/gs/lexer"
	"github.com/suzuken/gs/types"
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	lex := lexer.New()
	r := strings.NewReader("(print 1)")
	lex.Init(r)
	lex.Scan()
	parser := New(lex)
	actual := []types.Expression{
		"print",
		"1",
	}
	exps, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser failed: %s", err)
	}
	if !reflect.DeepEqual(exps, actual) {
		t.Fatalf("expressions is not expected. %v", exps)
	}
}

func TestParserRecursive(t *testing.T) {
	lex := lexer.New()
	r := strings.NewReader("(define (square x) (* x x))")
	lex.Init(r)
	lex.Scan()

	parser := New(lex)
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
	exps, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser failed: %s", err)
	}
	if !reflect.DeepEqual(exps, expected) {
		t.Fatalf("expressions is not expected. %v", exps)
	}
}
