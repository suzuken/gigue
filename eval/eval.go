// Package eval implements an evaluator of gigue.
package eval

import (
	"errors"
	"fmt"
	"github.com/suzuken/gigue/lexer"
	"github.com/suzuken/gigue/parser"
	"github.com/suzuken/gigue/types"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// evalDefine evaluate (define ...) style expression.
func evalDefine(exps []types.Expression, env *Env) (types.Expression, error) {
	if len(exps) < 2 {
		return nil, errors.New("define clause must have symbol and body")
	}
	switch tt := exps[1].(type) {
	// put symbol and variables
	// (define x 1) style definition
	case types.Symbol:
		value, err := Eval(exps[2], env)
		if err != nil {
			return nil, err
		}
		env.Put(tt, value)
		return nil, nil
	// (define (hoge args) (..)) style definition
	// above style is syntax sugar for lambda.
	case []types.Expression:
		if len(tt) < 2 {
			return nil, errors.New("define statament must have more than 2 words")
		}
		caddr, ok := tt[0].(types.Symbol)
		if !ok {
			return nil, errors.New("(define x) of x should be symbol")
		}
		// create lambda and put it into environment
		env.Put(caddr, Lambda{tt[1:], exps[2], env})
		return nil, nil
	default:
		return nil, nil
	}
}

// evalIf evaluates if-clause.
// predicate consequent alternative
// like, (if (ok? yeah) (go) (not go))
func evalIf(predicate, consequent, alternative types.Expression, env *Env) (types.Expression, error) {
	bb, err := evalPredicate(predicate, env)
	if err != nil {
		return nil, err
	}
	if bb {
		return Eval(consequent, env)
	}
	return Eval(alternative, env)
}

// evalPredicate is helper for evaluate expression should be boolean.
func evalPredicate(exp types.Expression, env *Env) (types.Boolean, error) {
	b, err := Eval(exp, env)
	if err != nil {
		return false, err
	}
	bb, ok := b.(types.Boolean)
	if !ok {
		return false, fmt.Errorf("the expression should return types.Boolean, exps: %v", exp)
	}
	return bb, nil
}

func evalCond(exps []types.Expression, env *Env) (types.Expression, error) {
	for _, operand := range exps[1:] {
		tt, ok := operand.([]types.Expression)
		if !ok {
			return nil, errors.New("cond clause must have expression")
		}
		if tt[0] == types.Symbol("else") {
			return Eval(tt[1], env)
		}
		if bb, err := evalPredicate(tt[0], env); err != nil {
			return nil, err
		} else if bb {
			return Eval(tt[1], env)
		}
	}
	// unreachable
	return nil, nil
}

func evalBegin(env *Env, exps ...types.Expression) (types.Expression, error) {
	var lastExp types.Expression
	for _, beginExp := range exps {
		l, err := Eval(beginExp, env)
		if err != nil {
			return nil, err
		}
		lastExp = l
	}
	return lastExp, nil
}

// evalLoad evaluates (load "file.scm") style definition.
// loading file and evaluate it.
func evalLoad(path string, env *Env) (types.Expression, error) {
	current, err := env.Get("#current-load-path")
	if err != nil {
		return nil, err
	}
	// if path is set, search from current directory
	if p := fmt.Sprintf("%s", current); p != "" {
		// if path start with /, deal as absolute path
		// if not, deal as relative path
		if !strings.HasPrefix(path, "/") {
			path = filepath.Join(filepath.Dir(p), path)
		}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// evaluate file in given environment.
	return EvalFile(abs, env)
}

func evalApplication(env *Env, operator types.Expression, operands ...types.Expression) (types.Expression, error) {
	// extend environment
	exps := make([]types.Expression, 0, len(operands)-1)
	for _, operand := range operands {
		exp, err := Eval(operand, env)
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
	}
	// maybe, it is primitive procedure or compound procedure.
	fn, err := Eval(operator, env)
	if err != nil {
		return nil, err
	}
	return Apply(fn, exps)
}

// Eval is body of evaluator
func Eval(exp types.Expression, env *Env) (types.Expression, error) {
	switch t := exp.(type) {
	case types.Boolean, types.Number, *types.Pair, string:
		return t, nil
	case types.Symbol:
		// it's variable or expression. get value from environment
		e, err := env.Get(t)
		if err != nil {
			return nil, err
		}
		return e, nil
	case []types.Expression:
		if len(t) == 0 {
			return &types.Pair{}, nil
		}
		switch t[0] {
		case types.Symbol("define"):
			return evalDefine(t, env)
		case types.Symbol("if"):
			if len(t) < 4 {
				return nil, errors.New("syntax error: if clause must be (if predicate consequent alternative) style")
			}
			return evalIf(t[1], t[2], t[3], env)
		case types.Symbol("cond"):
			if len(t) < 2 {
				return nil, errors.New("syntax error: cond clause must be (cond predicate consequent alternative) style")
			}
			return evalCond(t, env)
		case types.Symbol("lambda"):
			if len(t) < 3 {
				return nil, errors.New("lambda must have more than 3 words")
			}
			return Lambda{t[1], t[2], env}, nil
		case types.Symbol("begin"):
			return evalBegin(env, t[1:]...)
		case types.Symbol("load"):
			path, ok := t[1].(string)
			if !ok {
				return nil, errors.New("syntax error: args of load should be string")
			}
			return evalLoad(path, env)
		default:
			return evalApplication(env, t[0], t[1:]...)
		}
	default:
		// not found any known operands. failed.
		return nil, fmt.Errorf("unknown expression type -- %v", exp)
	}
}

// EvalFile evaluate given file
func EvalFile(filename string, env *Env) (types.Expression, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	env.Put("#current-load-path", filename)
	defer f.Close()
	return EvalReader(f, env)
}

// EvalReader evaluate scheme program from io.Reader
func EvalReader(r io.Reader, env *Env) (types.Expression, error) {
	l := lexer.New(r)
	p := parser.New(l)
	if _, err := env.Get("#current-load-path"); err != nil {
		env.Put("#current-load-path", "")
	}
	var exps types.Expression
	for {
		tokens, err := p.Parse()
		if err != nil {
			return nil, err
		}
		// TODO should handle unknown token.
		if tokens == types.Symbol("") {
			break
		}
		exps, err = Eval(tokens, env)
		if err != nil {
			return nil, err
		}
	}
	return exps, nil
}

// Apply receives procedure and arguments. if procedure is compounded, evaluate on extended environment.
func Apply(procedure types.Expression, args []types.Expression) (types.Expression, error) {
	switch p := procedure.(type) {
	case Lambda:
		// extend environment base on lambda arguments
		// should bind argument to this environment.
		// for example, (define (sum x y) (+ x y)) and given (sum 1 2),
		// then creates frames which have x = 1 and y = 2.
		env := &Env{m: make(Frame), parent: p.Env}
		env.Setup()
		switch lambdaArgs := p.Args.(type) {
		case []types.Expression:
			if len(lambdaArgs) != len(args) {
				return nil, errors.New("given args is not match with lambda args")
			}
			for i, arg := range lambdaArgs {
				env.Put(arg.(types.Symbol), args[i])
			}
		default:
			env.Put(lambdaArgs.(types.Symbol), lambdaArgs)
		}
		return Eval(p.Body, env)
	case func(...types.Expression) (types.Expression, error):
		// primitive procedure
		return p(args...)
	default:
		return nil, fmt.Errorf("unknown procedure type -- %v", procedure)
	}
}
