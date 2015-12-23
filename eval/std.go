package eval

import (
	"io"
	"strings"
)

// StandardLibrary is pseudo function returns standard library code of scheme.
// Easy for building.
func StandardLibrary() io.Reader {
	return strings.NewReader(`
(define (caar given)
  (car (car given)))

(define (caaar given)
  (car (caar given)))

(define (cadr given)
  (car (cdr given)))

(define (cddr given)
  (cdr (cdr given)))

(define (cdadr given)
  (cdr (car (cdr given))))

(define (cddar given)
  (cdr (cdr (car given))))

(define (cdddr given)
  (cdr (cdr (cdr given))))

(define (caddr given)
  (car (cdr (cdr given))))
`)
}
