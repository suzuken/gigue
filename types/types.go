package types

import (
	"fmt"
)

// S is S-expression
type Expression interface{}

// Number is number of scheme. (based on Go float64)
type Number float64

// Symbol is index of S-expression in some environment.
type Symbol string

// Boolean is boolean of scheme.
type Boolean bool

// String is string of scheme.
type String string

func (b Boolean) String() string {
	if b {
		return "#t"
	}
	return "#f"
}

// Pair is cons
type Pair struct {
	Car Expression
	Cdr Expression
}

func (p *Pair) String() string {
	return fmt.Sprintf("(%s . %s)", p.Car, p.Cdr)
}

func (p *Pair) IsNull() bool {
	return p.Car == nil && p.Cdr == nil
}

func (p *Pair) IsPair() bool {
	return !p.IsNull()
}

// IsList returns if pair is list or not.
//
// * empty pair is list
// * end of list should be empty pair (empty list)
func (p *Pair) IsList() bool {
	pp := p
	for {
		if pp.IsNull() {
			return true
		}
		switch cdr := pp.Cdr.(type) {
		case *Pair:
			pp = cdr
		default:
			return false
		}
	}
}

// Append add cons to given pair
func (p *Pair) Append(exp Expression) *Pair {
	// append exp to tail
	pp := p
	for {
		if pp.IsNull() {
			break
		}
		pp = pp.Cdr.(*Pair)
	}
	pp.Car = exp
	pp.Cdr = &Pair{}
	return pp
}

func NewList(args ...Expression) *Pair {
	p := &Pair{Car: nil, Cdr: nil}
	for _, arg := range args {
		p.Append(arg)
	}
	return p
}
