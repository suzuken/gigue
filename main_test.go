package main

import (
	"github.com/suzuken/gs/eval"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEv(t *testing.T) {
	env := eval.NewEnv()
	env.Setup()

	r := strings.NewReader(`
(define x 1)
(print x)
	`)
	if _, err := ev(r, env); err != nil {
		t.Fatalf("eval error read from io.Reader: %s", err)
	}
}

// visit generate WalkFunc for traversing examples directory.
func visit(t *testing.T) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		// skip directory
		if f.IsDir() {
			return nil
		}
		env := eval.NewEnv()
		env.Setup()
		if _, err := evalFile(path, env); err != nil {
			t.Fatalf("eval file failed. file: %s, err: %s\nenv: %#v", path, err, env)
			return err
		}
		t.Logf("eval file success: file: %s", path)
		return nil
	}
}

// TestExecute execute scheme scripts under examples in actual.
// simply check if cause error or not.
func TestExecute(t *testing.T) {
	if err := filepath.Walk("examples", visit(t)); err != nil {
		t.Fatalf("eval file failed: %s", err)
	}
}
