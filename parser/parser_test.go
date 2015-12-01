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
		types.Symbol("print"),
		types.Number(1),
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
		types.Symbol("define"),
		[]types.Expression{
			types.Symbol("square"),
			types.Symbol("x"),
		},
		[]types.Expression{
			types.Symbol("*"),
			types.Symbol("x"),
			types.Symbol("x"),
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

func TestParseBoolean(t *testing.T) {
	lex := lexer.New()
	r := strings.NewReader("(print #t)")
	lex.Init(r)
	lex.Scan()
	parser := New(lex)
	actual := []types.Expression{
		types.Symbol("print"),
		types.Boolean(true),
	}
	exps, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser failed: %s", err)
	}
	if !reflect.DeepEqual(exps, actual) {
		t.Fatalf("expressions is not expected. %v", exps)
	}
}
