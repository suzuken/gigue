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

	if _, err := Eval(exp, env); err != nil {
		t.Fatalf("eval but error : %s", err)
	}
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

func TestEvalDefineAndCaliculation(t *testing.T) {
	env := NewEnv()
	env.Setup()

	if _, err := Eval([]types.Expression{
		types.Symbol("define"),
		types.Symbol("x"),
		types.Number(1),
	}, env); err != nil {
		t.Fatalf("eval but error : %s", err)
	}

	// (+ x 3 (* 5 (- 5 2)))
	exp := []types.Expression{
		types.Symbol("+"),
		types.Symbol("x"),
		types.Number(3),
		[]types.Expression{
			types.Symbol("*"),
			types.Number(5),
			[]types.Expression{
				types.Symbol("-"),
				types.Number(5),
				types.Number(2),
			},
		},
	}

	r, err := Eval(exp, env)
	if err != nil {
		t.Fatalf("eval but error : %s", err)
	}

	if r != types.Number(19) {
		t.Fatalf("given x = 1, (+ x 3 (* 5 (- 5 2))) should equal 19 but get: %v", r)
	}
}
