package eval

import (
	"github.com/suzuken/gs/types"
	"testing"
)

func TestEval(t *testing.T) {
	env := NewEnv()
	env.Setup()

	exp := []types.Expression{
		types.Symbol("print"),
		types.Number(1),
	}

	r, err := Eval(exp, env)
	if err != nil {
		t.Fatalf("eval but error : %s", err)
	}
	t.Logf("%v\n", r)
}

func TestEvalSum(t *testing.T) {
	env := NewEnv()
	env.Setup()

	exp := []types.Expression{
		types.Symbol("+"),
		types.Number(1),
		types.Number(2),
	}

	r, err := Eval(exp, env)
	if err != nil {
		t.Fatalf("eval but error : %s", err)
	}
	if r != types.Number(3) {
		t.Fatalf("1 + 2 should 3 but not: %v", r)
	}

}
