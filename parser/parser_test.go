package parser

import (
	"github.com/suzuken/gigue/lexer"
	"github.com/suzuken/gigue/types"
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

func TestParseString(t *testing.T) {
	lex := lexer.New()
	r := strings.NewReader("(print \"it's test\")")
	lex.Init(r)
	lex.Scan()
	parser := New(lex)
	actual := []types.Expression{
		types.Symbol("print"),
		"it's test",
	}
	exps, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser failed: %s", err)
	}
	if !reflect.DeepEqual(exps, actual) {
		t.Fatalf("expressions is not expected. %#v", exps)
	}
}

func TestParseDash(t *testing.T) {
	lex := lexer.New()
	r := strings.NewReader("(define a-b-c-efg x)")
	lex.Init(r)
	lex.Scan()
	parser := New(lex)
	actual := []types.Expression{
		types.Symbol("define"),
		types.Symbol("a-b-c-efg"),
		types.Symbol("x"),
	}
	exps, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser failed: %s", err)
	}
	if !reflect.DeepEqual(exps, actual) {
		t.Fatalf("expressions is not expected. %#v", exps)
	}
}

func TestParseLineDelimited(t *testing.T) {
	lex := lexer.New()
	r := strings.NewReader(`
(define (fib n)
  (cond ((= n 0) 0)
        ((= n 1) 1)
        (else (+ (fib (- n 1)) (fib (- n 2))))))
	`)
	lex.Init(r)
	lex.Scan()
	parser := New(lex)
	actual := []types.Expression{
		types.Symbol("define"),
		[]types.Expression{
			types.Symbol("fib"),
			types.Symbol("n"),
		},
		[]types.Expression{
			types.Symbol("cond"),
			[]types.Expression{
				[]types.Expression{
					types.Symbol("="),
					types.Symbol("n"),
					types.Number(0),
				},
				types.Number(0),
			},
			[]types.Expression{
				[]types.Expression{
					types.Symbol("="),
					types.Symbol("n"),
					types.Number(1),
				},
				types.Number(1),
			},
			[]types.Expression{
				types.Symbol("else"),
				[]types.Expression{
					types.Symbol("+"),
					[]types.Expression{
						types.Symbol("fib"),
						[]types.Expression{
							types.Symbol("-"),
							types.Symbol("n"),
							types.Number(1),
						},
					},
					[]types.Expression{
						types.Symbol("fib"),
						[]types.Expression{
							types.Symbol("-"),
							types.Symbol("n"),
							types.Number(2),
						},
					},
				},
			},
		},
	}
	exps, err := parser.Parse()
	if err != nil {
		t.Fatalf("parser failed: %s", err)
	}
	if !reflect.DeepEqual(exps, actual) {
		t.Fatalf("expressions is not expected. %#v", exps)
	}
}
