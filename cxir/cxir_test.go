package cxir

import (
	"strings"
	"testing"

	"github.com/skycoin/skycoin-cx-forth/forth"
)

func TestExamples(t *testing.T) {
	for _, ex := range Examples {
		t.Run(ex.Name, func(t *testing.T) {
			var b strings.Builder
			if err := Run(&ex.Program, &b); err != nil {
				src, _ := Translate(&ex.Program)
				t.Fatalf("%v\n%s", err, src)
			}
			if b.String() != ex.Output {
				t.Errorf("got %q want %q", b.String(), ex.Output)
			}
		})
	}
}

func TestTranslateErrors(t *testing.T) {
	main := func(e ...Expr) *Program {
		return &Program{Functions: []Function{{Name: "main", Locals: []Arg{Var("x", I32)}, Exprs: e}}}
	}
	cases := []struct {
		p    *Program
		want string
	}{
		{&Program{}, "no main"},
		{main(Print(Var("y", I32))), "undefined variable y"},
		{main(Goto(5)), "jump target out of range"},
		{main(Call("nope", nil)), "undefined function nope"},
		{main(Op1(OpAdd, Int32(1), Int32(1), Int32(1))), "assign to a literal"},
		{main(Print(Str(`a"b`))), "quote"},
		{&Program{Functions: []Function{
			{Name: "main", Exprs: []Expr{Call("a", nil)}},
			{Name: "a", Exprs: []Expr{Call("b", nil)}},
			{Name: "b", Exprs: []Expr{Call("a", nil)}},
		}}, "mutual recursion"},
	}
	for _, c := range cases {
		_, err := Translate(c.p)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("got %v want %q", err, c.want)
		}
	}
}

// A runtime error inside translated code surfaces as a VM error.
func TestRuntimeError(t *testing.T) {
	p := &Program{Functions: []Function{{
		Name:   "main",
		Locals: []Arg{Var("x", I32)},
		Exprs:  []Expr{Op1(OpDiv, Var("x", I32), Int32(1), Int32(0))},
	}}}
	err := Run(p, nil)
	if _, ok := err.(*forth.Error); !ok || !strings.Contains(err.Error(), "division by zero") {
		t.Fatalf("got %v", err)
	}
}
