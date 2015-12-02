package eval

import (
	"errors"
	"github.com/suzuken/gs/types"
)

// Eval is body of evaluator
func Eval(exp types.Expression, env *Env) (types.Expression, error) {
	switch t := exp.(type) {
	case types.Boolean:
		return t, nil
	case types.Number:
		return t, nil
	case types.Symbol:
		// it's variable. get value from environment
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
				return tt, nil
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
				env.Put(caddr, types.Lambda{tt[0], t[2]})
				return tt[0], nil
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
			// TODO transform and use if evaluation here, too.
		case "lambda":
			if len(t) < 3 {
				return nil, errors.New("lambda must have more than 3 words.")
			}
			return types.Lambda{t[1], t[2]}, nil
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

// Apply receives procedure and arguments. if procedure is compounded, evaluate on extended environment.
func Apply(procedure types.Expression, args []types.Expression) (types.Expression, error) {
	switch p := procedure.(type) {
	case func(...types.Expression) (types.Expression, error):
		// primitive procedure
		return p(args...)
	default:
		return nil, errors.New("Unknown procedure type")
	}
	return nil, nil
}

// EvalSeauencd evaluate sequence of expressions in certain environment.
// Return is last expression.
func EvalSequence(exps []types.Expression, env *Env) (types.Expression, error) {
	if len(exps) == 1 {
		return Eval(exps[0], env)
	}
	// making environment (Yes, it's pointer)
	if _, err := Eval(exps[0], env); err != nil {
		return nil, err
	}
	return EvalSequence(exps[1:], env)
}

// listOfValues returns arguments for evaluator.
func listOfValues(exps []types.Expression, env *Env) (types.Expression, error) {
	if len(exps) <= 0 {
		return nil, nil
	}
	// evaluate exps one by one on each environment
	first, err := Eval(exps[0], env)
	if err != nil {
		return nil, err
	}
	// TODO: should use for in Go way?
	rest, err := listOfValues(exps[1:], env)
	if err != nil {
		return nil, err
	}
	return types.Pair{first, rest}, nil
}
