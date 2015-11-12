package lexer

import (
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	lex := New()
	r := strings.NewReader("(print 1)")
	lex.Init(r)
	lex.NextToken()
	tokens, err := lex.Scan()
	if err != nil {
		t.Fatalf("lexer failed: %s", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("tokens is not expected. %v", tokens)
	}
	if tokens[0] != "print" || tokens[1] != "1" {
		t.Fatalf("tokens is not expected. %v", tokens)
	}
}
