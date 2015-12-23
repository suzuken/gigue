package eval

import (
	"errors"
	"fmt"
	"github.com/suzuken/gigue/types"
	"sync"
)

// Frame is symbol to Expression map
type Frame map[types.Symbol]types.Expression

// Env is scheme environment for evaluation
type Env struct {
	sync.RWMutex
	m      Frame // m is symbol table for expression
	parent *Env  // parent is parent Environment. Env is nested.
}

// NewEnv creates new environment
func NewEnv() *Env {
	symbols := make(Frame) // TODO: more flexible stack size control
	return &Env{m: symbols, parent: nil}
}

// Extend extends environments by given variables and values.
func (e *Env) Extend(frame Frame) {
	for s, exp := range frame {
		e.Put(s, exp)
	}
}

// Setup returns initial environment for evaluation.
// Primitive procedures is defined by Go, inner intepriter.
// Standard library is written in Scheme, extend it's functionality.
func (e *Env) Setup() {
	e.Extend(NewPrimitiveProcedureFrame())
	e.LoadStandardLibrary()
}

// LoadStandardLibrary load gigue's standard library and create it's environment.
func (env *Env) LoadStandardLibrary() error {
	if _, err := EvalReader(StandardLibrary(), env); err != nil {
		return err
	}
	return nil
}

// NewPrimitiveProcedureFrame returns frame of primitive procedures.
// Gigue provides functionality of base scheme procedure.
func NewPrimitiveProcedureFrame() Frame {
	return Frame{
		"car":     Car,
		"cdr":     Cdr,
		"cons":    Cons,
		"print":   Print,
		"+":       Add,
		"-":       Subtract,
		"*":       Multiply,
		"/":       Divide,
		">":       GreaterThan,
		"<":       LessThan,
		">=":      GreaterThanEqual,
		"<=":      LessThanEqual,
		"eq?":     IsEqual,
		"=":       IsEqual,
		"null?":   IsNull,
		"list":    List,
		"list?":   IsList,
		"string?": IsString,
		"symbol?": IsSymbol,
	}
}

// Car is implementation of car
func Car(args ...types.Expression) (types.Expression, error) {
	a, ok := args[0].(*types.Pair)
	if !ok {
		return nil, errors.New("arguments of car should pair")
	}
	return a.Car, nil
}

// Cdr is cdr
func Cdr(args ...types.Expression) (types.Expression, error) {
	a, ok := args[0].(*types.Pair)
	if !ok {
		return nil, errors.New("arguments of cdr should pair")
	}
	return a.Cdr, nil
}

func Cons(args ...types.Expression) (types.Expression, error) {
	return &types.Pair{Car: args[0], Cdr: args[1]}, nil
}

func Print(args ...types.Expression) (types.Expression, error) {
	if len(args) == 1 {
		fmt.Println(args[0])
	} else {
		fmt.Println(args)
	}
	return nil, nil
}

func Add(args ...types.Expression) (types.Expression, error) {
	sum, ok := args[0].(types.Number)
	if !ok {
		return nil, fmt.Errorf("given args is not number: %#v", args[0])
	}
	for _, adder := range args[1:] {
		sum = sum + adder.(types.Number)
	}
	return sum, nil
}

func Subtract(args ...types.Expression) (types.Expression, error) {
	sub, ok := args[0].(types.Number)
	if !ok {
		return nil, fmt.Errorf("given args is not number: %v", args[0])
	}
	for _, s := range args[1:] {
		sub = sub - s.(types.Number)
	}
	return sub, nil
}

func Multiply(args ...types.Expression) (types.Expression, error) {
	mul, ok := args[0].(types.Number)
	if !ok {
		return nil, fmt.Errorf("given args is not number: %v", args[0])
	}
	for _, m := range args[1:] {
		mul = mul * m.(types.Number)
	}
	return mul, nil
}

func Divide(args ...types.Expression) (types.Expression, error) {
	div, ok := args[0].(types.Number)
	if !ok {
		return nil, fmt.Errorf("given args is not number: %v", args[0])
	}
	for _, d := range args[1:] {
		div = div / d.(types.Number)
	}
	return div, nil
}

func GreaterThan(args ...types.Expression) (types.Expression, error) {
	return types.Boolean(args[0].(types.Number) > args[1].(types.Number)), nil
}

func GreaterThanEqual(args ...types.Expression) (types.Expression, error) {
	return types.Boolean(args[0].(types.Number) >= args[1].(types.Number)), nil
}

func LessThan(args ...types.Expression) (types.Expression, error) {
	return types.Boolean(args[0].(types.Number) < args[1].(types.Number)), nil
}

func LessThanEqual(args ...types.Expression) (types.Expression, error) {
	return types.Boolean(args[0].(types.Number) <= args[1].(types.Number)), nil
}

func IsEqual(args ...types.Expression) (types.Expression, error) {
	return types.Boolean(args[0] == args[1]), nil
}

func IsNull(args ...types.Expression) (types.Expression, error) {
	pair, ok := args[0].(*types.Pair)
	if !ok {
		return types.Boolean(false), nil
	}
	return types.Boolean(pair.IsNull()), nil
}

func List(args ...types.Expression) (types.Expression, error) {
	return types.NewList(args...), nil
}

func IsList(args ...types.Expression) (types.Expression, error) {
	pair, ok := args[0].(*types.Pair)
	if !ok {
		return types.Boolean(false), nil
	}
	return types.Boolean(pair.IsList()), nil
}

func IsString(args ...types.Expression) (types.Expression, error) {
	if _, ok := args[0].(string); !ok {
		return types.Boolean(false), nil
	}
	return types.Boolean(true), nil
}

func IsSymbol(args ...types.Expression) (types.Expression, error) {
	if _, ok := args[0].(types.Symbol); !ok {
		return types.Boolean(false), nil
	}
	return types.Boolean(true), nil
}

// Put creates new symbol to table
func (e *Env) Put(s types.Symbol, exp types.Expression) {
	e.Lock()
	defer e.Unlock()
	e.m[s] = exp
}

// Get fetch expression by symbol from environment
func (e *Env) Get(s types.Symbol) (types.Expression, error) {
	e.RLock()
	defer e.RUnlock()
	v, ok := e.m[s]
	if !ok {
		if e.parent != nil {
			return e.parent.Get(s)
		}
		// if symbol not found, return itself
		return s, nil
	}
	return v, nil
}

// Remove symbol from environment
func (e *Env) Remove(s types.Symbol) {
	e.Lock()
	defer e.Unlock()
	delete(e.m, s)
}

// Lambda is definition of lambda
type Lambda struct {
	// Args are temporary parameters
	Args types.Expression
	// Body is expression to evalute
	Body types.Expression
	// Env is environent for evaluate this lambda function
	Env *Env
}
