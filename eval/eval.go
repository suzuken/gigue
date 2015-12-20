package eval

import (
	"errors"
	"github.com/suzuken/gs/lexer"
	"github.com/suzuken/gs/parser"
	"github.com/suzuken/gs/types"
	"io"
	"os"
	"path/filepath"
	"text/scanner"
)

// Eval is body of evaluator
func Eval(exp types.Expression, env *Env) (types.Expression, error) {
	switch t := exp.(type) {
	case types.Boolean, types.Number, types.String:
		return t, nil
	case types.Symbol:
		// it's variable or expression. get value from environment
		e, err := env.Get(t)
		if err != nil {
			return nil, err
		}
		return e, nil
	case []types.Expression:
		// this is multiple expressions pattern
		// at first, get car. car of expression is symbol for each expression
		car, ok := t[0].(types.Symbol)
		if !ok {
			return nil, errors.New("cannot conversion car of expressions. it should be types.Symbol but not.")
		}
		switch car {
		case "define":
			if len(t) < 2 {
				return nil, errors.New("define clause must have symbol and body.")
			}
			// TODO consider if simply Eval(t[2], env) and put it into env.
			// because []types.Expression should be evaluated by Eval() anyway.
			switch tt := t[1].(type) {
			// put symbol and variables
			// (define x 1) style definition
			case types.Symbol:
				value, err := Eval(t[2], env)
				if err != nil {
					return nil, err
				}
				env.Put(tt, value)
				return nil, nil
			// (define (hoge args) (..)) style definition
			case []types.Expression:
				if len(tt) < 2 {
					return nil, errors.New("define statament must have more than 2 words.")
				}
				caddr, ok := tt[0].(types.Symbol)
				if !ok {
					return nil, errors.New("(define x) of x should be symbol..")
				}
				// create lambda and put it into environment
				// TODO verify if it's own func name is certainly put into environment.
				// it's necessary for evaluating recursive function.
				env.Put(caddr, Lambda{tt[1:], t[2], env})
				return nil, nil
			}
		case "if":
			// like, (if (ok? yeah) (go) (not go))
			b, err := Eval(t[1], env)
			if err != nil {
				return nil, err
			}
			bb, ok := b.(types.Boolean)
			if !ok {
				return nil, errors.New("if-predicate should return types.Boolean")
			}
			if bb {
				return Eval(t[2], env)
			} else {
				return Eval(t[3], env)
			}
		case "cond":
			for _, operand := range t[1:] {
				tt, ok := operand.([]types.Expression)
				if !ok {
					return nil, errors.New("cond clause must have expression")
				}
				if tt[0] == types.Symbol("else") {
					return Eval(tt[1], env)
				}
				b, err := Eval(tt[0], env)
				if err != nil {
					return nil, err
				}
				bb, ok := b.(types.Boolean)
				if ok && bb == true {
					return Eval(tt[1], env)
				}
			}
		case "lambda":
			if len(t) < 3 {
				return nil, errors.New("lambda must have more than 3 words.")
			}
			return Lambda{t[1], t[2], env}, nil
		case "begin":
			// (begin s1 s2 ... last)
			var lastExp types.Expression
			for _, beginExp := range t[1:] {
				l, err := Eval(beginExp, env)
				if err != nil {
					return nil, err
				}
				lastExp = l
			}
			return lastExp, nil
		case "load":
			// (load "file.scm") loading file and evaluate it.
			path, ok := t[1].(string)
			if !ok {
				return nil, errors.New("args of load should be string.")
			}
			abs, err := filepath.Abs(path)
			if err != nil {
				return nil, err
			}
			// evaluate file in given environment.
			return EvalFile(abs, env)
		default:
			// extend environment
			exps := make([]types.Expression, 0)
			for _, operand := range t[1:] {
				exp, err := Eval(operand, env)
				if err != nil {
					return nil, err
				}
				exps = append(exps, exp)
			}
			// maybe, it is primitive procedure or compound procedure.
			fn, err := Eval(car, env)
			if err != nil {
				return nil, err
			}
			result, err := Apply(fn, exps)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	default:
		// not found any known operands. failed.
		return nil, errors.New("unknown expression type")
	}
	return nil, nil
}

// EvalFile evaluate given file
func EvalFile(filename string, env *Env) (types.Expression, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return EvalReader(f, env)
}

// ev evaluate scheme program from io.Reader
func EvalReader(r io.Reader, env *Env) (types.Expression, error) {
	l := lexer.New()
	l.Init(r)
	l.Scan()
	p := parser.New(l)
	for l.Token != scanner.EOF {
		exps, err := p.Parse()
		if err != nil {
			return nil, err
		}
		if _, err := Eval(exps, env); err != nil {
			return nil, err
		}
		l.Scan()
	}
	return nil, nil
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
		return nil, errors.New("Unknown procedure type")
	}
}
