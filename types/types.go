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
