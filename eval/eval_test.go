package eval

import (
	"github.com/suzuken/gigue/types"
	"os"
	"path/filepath"
	"strings"
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

func TestEvalDefineFunction(t *testing.T) {
	env := NewEnv()
	env.Setup()

	// (define (sum x y) (+ x y))
	if _, err := Eval([]types.Expression{
		types.Symbol("define"),
		[]types.Expression{
			types.Symbol("sum"),
			types.Symbol("x"),
			types.Symbol("y"),
		},
		[]types.Expression{
			types.Symbol("+"),
			types.Symbol("x"),
			types.Symbol("y"),
		},
	}, env); err != nil {
		t.Fatalf("eval but error : %s", err)
	}

	// (sum 1 2)
	exp := []types.Expression{
		types.Symbol("sum"),
		types.Number(1),
		types.Number(2),
	}

	r, err := Eval(exp, env)
	if err != nil {
		t.Fatalf("eval but error : %s", err)
	}

	if r != types.Number(3) {
		t.Fatalf("given (sum 1 2) should equal 3 but get: %v", r)
	}
}

func TestEvalRecursiveFunction(t *testing.T) {
	env := NewEnv()
	env.Setup()

	// (define (fib n)
	//   (cond ((= n 0) 0)
	//         ((= n 1) 1)
	//         (else (+ (fib (- n 1)) (fib (- n 2))))))
	if _, err := Eval([]types.Expression{
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
	}, env); err != nil {
		t.Fatalf("eval but error : %s", err)
	}

	exp := []types.Expression{
		types.Symbol("fib"),
		types.Number(10),
	}

	r, err := Eval(exp, env)
	if err != nil {
		t.Fatalf("eval but error : %s, env: %v", err, env)
	}

	if r != types.Number(55) {
		t.Fatalf("given (fib 10) should equal 55 but get: %v", r)
	}
}

func TestPrimitiveListOperation(t *testing.T) {
	env := NewEnv()
	env.Setup()

	// (define x (list 1 2 3))
	if _, err := Eval([]types.Expression{
		types.Symbol("define"),
		types.Symbol("x"),
		[]types.Expression{
			types.Symbol("list"),
			types.Number(1),
			types.Number(2),
			types.Number(3),
		},
	}, env); err != nil {
		t.Fatalf("eval error: %s", err)
	}

	car, err := Eval([]types.Expression{
		types.Symbol("car"),
		types.Symbol("x"),
	}, env)
	if err != nil {
		t.Fatalf("eval error: %s", err)
	}
	if car != types.Number(1) {
		t.Fatal("cannot get car")
	}

	cdrExp, err := Eval([]types.Expression{
		types.Symbol("cdr"),
		types.Symbol("x"),
	}, env)
	if err != nil {
		t.Fatalf("eval error: %s", err)
	}
	cdr, ok := cdrExp.(*types.Pair)
	if !ok {
		t.Fatal("cdr should be pair but not")
	}
	if cdr.Car != types.Number(2) {
		t.Fatal("cannot get cdr")
	}
}

func TestEvalReader(t *testing.T) {
	env := NewEnv()
	env.Setup()

	r := strings.NewReader(`
(define x 1)
(print x)
	`)
	if _, err := EvalReader(r, env); err != nil {
		t.Fatalf("eval error read from io.Reader: %s", err)
	}
}

// visit generate WalkFunc for traversing examples directory.
func visit(t *testing.T) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		// skip directory
		if f.IsDir() {
			return nil
		}
		env := NewEnv()
		env.Setup()
		if _, err := EvalFile(path, env); err != nil {
			t.Fatalf("eval file failed. file: %s, err: %s\nenv: %#v", path, err, env)
			return err
		}
		t.Logf("eval file success: file: %s", path)
		return nil
	}
}

// TestExecute execute scheme scripts under examples in actual.
// simply check if cause error or not.
func TestExecute(t *testing.T) {
	if err := filepath.Walk("../examples", visit(t)); err != nil {
		t.Fatalf("eval file failed: %s", err)
	}
}
