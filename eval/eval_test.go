package eval

import (
	"github.com/suzuken/gs/lexer"
	"testing"
)

func TestEval(t *testing.T) {
	Eval(lexer.New(), nil)
}
