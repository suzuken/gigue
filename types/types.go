package types

import (
	"fmt"
	"sync"
)

// S is S-expression
type Expression interface{}

// Number is number of scheme. (based on Go float64)
type Number float64

// Symbol is index of S-expression in some environment.
type Symbol string

// Boolean is boolean of scheme.
type Boolean bool

func (b Boolean) String() string {
	if b {
		return "#t"
	}
	return "#f"
}

// Env is scheme environment for evaluation
type Env struct {
	*sync.RWMutex
	m      map[Symbol]*Expression // m is symbol table for expression
	parent *Env                   // parent is parent Environment. Env is nested.
}

// NewEnv creates new environment
func NewEnv() *Env {
	symbols := make(map[Symbol]*Expression) // TODO: more flexible stack size control
	return &Env{m: symbols}
}

// Put creates new symbol to table
func (e *Env) Put(s Symbol, exp *Expression) {
	e.Lock()
	defer e.Unlock()
	e.m[s] = exp
}

// Get fetch expression by symbol from environment
func (e *Env) Get(s Symbol) (*Expression, error) {
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
func (e *Env) Remove(s Symbol) {
	e.Lock()
	defer e.Unlock()
	delete(e.m, s)
}

// Pair is cons
type Pair struct {
	Car Expression
	Cdr Expression
}

func (p *Pair) String() string {
	return fmt.Sprintf("(%s . %s)", p.Car, p.Cdr)
}

// List is list of scheme
type List struct {
	*Pair
}

func (l List) String() (str string) {
	// TODO implementation
	str = str + "("
	if l.Car == nil {
		str = str + ")"
		return str
	} else {
		// l.cdr is list
		str = str + fmt.Sprint(l.Cdr)
		return str
	}
	return str
}

// Len returns length of list
func (l *List) Len(num int) int {
	length := num
	list, ok := l.Cdr.(List)
	if !ok {
		// TODO should return error
		return length
	}
	return list.Len(length + 1)
}
