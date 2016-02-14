package types

import (
	"fmt"
	"strings"
)

// Expression is S-expression
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

// Pair is cons
type Pair struct {
	Car Expression
	Cdr Expression
}

func (p *Pair) String() string {
	if p.IsNull() {
		return "()"
	}
	if p.IsList() {
		var tokens []string
		pp := p
		for {
			if pp.IsNull() {
				break
			}
			tokens = append(tokens, fmt.Sprintf("%v", pp.Car))
			switch cdr := pp.Cdr.(type) {
			case *Pair:
				pp = cdr
			default:
				break
			}
		}
		return fmt.Sprintf("(%s)", strings.Join(tokens, " "))
	}
	return fmt.Sprintf("(%v . %v)", p.Car, p.Cdr)
}

// IsNull checking if pair is null or not.
func (p *Pair) IsNull() bool {
	return p.Car == nil && p.Cdr == nil
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

// Append add cons pair to given pair
// exp, first arguments of callee, should be car of this pair.
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

// NewList makes concatenated pair's list.
// Internally, creating last pair and concatenate it with previous pair.
func NewList(args ...Expression) *Pair {
	// In normal, p is prefer to be defined by var statement because of no allocation.
	// But in this case, for empty list should be return empty pair.
	p, prev := &Pair{}, &Pair{}
	for i := len(args) - 1; i >= 0; i-- {
		p = &Pair{args[i], prev}
		prev = p
	}
	return p
}
