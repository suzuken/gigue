# Gigue

[![Build Status](https://travis-ci.org/suzuken/gigue.svg)](https://travis-ci.org/suzuken/gigue)

One day after reading [SICP](https://mitpress.mit.edu/sicp/), I interested in implement my own scheme interpriter for practice. A design of Gigue is based on scheme interpriter of [SICP Chapter 4](https://mitpress.mit.edu/sicp/full-text/book/book-Z-H-25.html#%_chap_4).

Gigue is not [the metacircular evaluator](https://mitpress.mit.edu/sicp/full-text/book/book-Z-H-26.html#%_sec_4.1), simple scheme interpriter written in Go.

## Installation

    go get github.com/suzuken/gigue

## Examples

From [examples/fold.scm](examples/fold.scm),

```scheme
(define (fold-right op initial sequence)
  (if (null? sequence)
    initial
    (op (car sequence)
        (fold-right op initial (cdr sequence)))))

(print (fold-right / 1 (list 1 2 3)))
```

and `gigue examples/fold.scm`, then output `1.5`.

### REPL

Gigue has REPL.

```
-> % ./gigue
> (define x 1)
> (define y 2)
> (load "examples/sum.scm")
> (print (sum x y))
3
```

## Features

* `+`, `-`, `*`, `/`
* `cons`, `car`, `cdr`, `list`
* `<`, `>`
* `define`, `lambda`, `if`, `cond`, `begin`
* `load`

## LICENSE

MIT

## Author

Kenta Suzuki (a.k.a. suzuken)

## For your reference

Written in above, Gigue is based on the interpriter written in SICP Chapter 4. If you are not similar to scheme, you don't mind it. I think writing scheme interpriter in scheme - metacircular evaluator - is much simpler than written in Go or other language.

[SICP Chapter 4](https://mitpress.mit.edu/sicp/full-text/book/book-Z-H-25.html#%_chap_4)

If you're not have enough time to read SICP, [scm.go, a Scheme interpreter in Go, as in SICP and lis.py | De Babbelbox of Pieter Kelchtermans](https://pkelchte.wordpress.com/2013/12/31/scm-go/) may be helpful for understanding what a scheme interpriter is. [A minimal Scheme interpreter, as seen in lis.py and SICP](https://gist.github.com/pkelchte/c2bd76b9f8f9cd603b3c) is a minimal scheme interpriter (surprisingly only 250 lines) written in Go .

And I found some projects which implements scheme interpriter in Go.

* [k0kubun/gosick](https://github.com/k0kubun/gosick)
* [chrisbutcher/goscheme](https://github.com/chrisbutcher/goscheme)
* [chenyukang/GoScheme](https://github.com/chenyukang/GoScheme)
* [kedebug/LispEx](https://github.com/kedebug/LispEx)
* [bytbox/kakapo](https://github.com/bytbox/kakapo)
* [jlippitt/golisp](https://github.com/jlippitt/golisp)
