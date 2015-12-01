package eval

import (
	"fmt"
	"github.com/suzuken/gs/types"
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
func (e *Env) Setup() {
	e.Extend(NewPrimitiveProcedureFrame())
}

// NewPrimitiveProcedureFrame returns frame of primitive procedures.
// GS provides functionality of base scheme procedure.
func NewPrimitiveProcedureFrame() Frame {
	return Frame{
		"car":   Car,
		"cdr":   Cdr,
		"cons":  Cons,
		"print": Print,
		"+":     Add,
		"-":     Subtract,
		"*":     Multiply,
		"/":     Divide,
		">":     GreaterThan,
		"<":     LessThan,
		"eq?":   IsEqual,
		"null?": IsNull,
	}
}

// Car is implementation of car
func Car(args ...types.Expression) types.Expression {
	return args[0]
}

func Cdr(args ...types.Expression) types.Expression {
	return args[1:]
}

func Cons(args ...types.Expression) types.Expression {
	return types.Pair{args[0], args[1]}
}

func Print(args ...types.Expression) types.Expression {
	fmt.Print(args)
	return nil
}

func Add(args ...types.Expression) types.Expression {
	sum := args[0].(types.Number)
	for _, adder := range args[1:] {
		sum = sum + adder.(types.Number)
	}
	return sum
}

func Subtract(args ...types.Expression) types.Expression {
	sub := args[0].(types.Number)
	for _, s := range args[1:] {
		sub = sub - s.(types.Number)
	}
	return sub
}

func Multiply(args ...types.Expression) types.Expression {
	mul := args[0].(types.Number)
	for _, m := range args[1:] {
		mul = mul * m.(types.Number)
	}
	return mul
}

func Divide(args ...types.Expression) types.Expression {
	div := args[0].(types.Number)
	for _, d := range args[1:] {
		div = div / d.(types.Number)
	}
	return div
}

func GreaterThan(args ...types.Expression) types.Expression {
	return types.Boolean(args[0].(types.Number) > args[1].(types.Number))
}

func LessThan(args ...types.Expression) types.Expression {
	return types.Boolean(args[0].(types.Number) < args[1].(types.Number))
}

func IsEqual(args ...types.Expression) types.Expression {
	return types.Boolean(args[0] == args[1])
}

func IsNull(args ...types.Expression) types.Expression {
	return types.Boolean(args[0] == nil)
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
		// or return nil?
		return nil, fmt.Errorf("symbol not found from the environment: symbol %s", s)
	}
	return v, nil
}

// Remove symbol from environment
func (e *Env) Remove(s types.Symbol) {
	e.Lock()
	defer e.Unlock()
	delete(e.m, s)
}
