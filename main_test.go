package main

import (
	"github.com/suzuken/gs/eval"
	"strings"
	"testing"
)

func TestEv(t *testing.T) {
	env := eval.NewEnv()
	env.Setup()

	r := strings.NewReader(`
(define x 1)
(print x)
	`)
	if _, err := ev(r, env); err != nil {
		t.Fatalf("eval error read from io.Reader: %s", err)
	}
}
