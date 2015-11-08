package types

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

// S is S-expression
type S interface{}

// Number is number of scheme. (based on Go float64)
type Number float64

// Symbol is index of S-expression in some environment.
type Symbol string

// Boolean is boolean of scheme.
type Boolean bool

func (b *Boolean) String() string {
	if b {
		return "#t"
	}
	return "#f"
}

// Env is scheme environment for evaluation
type Env struct {
	*sync.RWMutex
	m      map[Symbol]*S // m is symbol table for expression
	parent *Env          // parent is parent Environment. Env is nested.
}

// NewEnv creates new environment
func NewEnv() *Env {
	symbols := make(map[Symbol]*S) // TODO: more flexible stack size control
	return &Env{m: symbols}
}

// Put creates new symbol to table
func (e *Env) Put(s Symbol, exp *S) {
	e.Lock()
	defer e.Unlock()
	e.m[s] = exp
}

// Get fetch expression by symbol from environment
func (e *Env) Get(s Symbol) (*S, error) {
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
	delete(e.m(s))
}
