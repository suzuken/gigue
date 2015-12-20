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

func verifyTokens(t *testing.T, x, y []string) bool {
	if len(x) != len(y) {
		t.Fatalf("token length is different: %v %v", x, y)
		return false
	}
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			t.Fatalf("#%d token is not same: %v %v", i, x, y)
			return false
		}
	}
	return true
}

// verifyTokenAll is useful method for testing given string how lexical analysis is works.
func verifyTokenAll(t *testing.T, given string, expected string) {
	lex := New()
	lex.Init(strings.NewReader(given))
	tokens, err := lex.TokenAll()
	if err != nil {
		t.Fatalf("lexer failed: %s", err)
	}
	verifyTokens(t, tokens, strings.Split(expected, ","))
}

func TestLexDefine(t *testing.T) {
	verifyTokenAll(t, "(print 1)", "(,print,1,)")
	verifyTokenAll(t, "(print \"test\")", "(,print,\"test\",)")
	verifyTokenAll(t, "(print\n1)\n\n", "(,print,1,),")
}
