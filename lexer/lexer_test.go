package lexer

import (
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	lex := New()
	r := strings.NewReader("(print 1)")
	lex.Init(r)
	lex.Scan()
	if lex.TokenText() != "(" {
		t.Fatalf("first string is ( but %s", lex.TokenText())
	}
	lex.Scan()
	if lex.TokenText() != "print" {
		t.Fatalf("second string is print but %s", lex.TokenText())
	}
	lex.Scan()
	if lex.TokenText() != "1" {
		t.Fatalf("third string is print but %s", lex.TokenText())
	}
	lex.Scan()
	if lex.TokenText() != ")" {
		t.Fatalf("forth string is print but %s", lex.TokenText())
	}
}
