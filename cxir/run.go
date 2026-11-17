package cxir

import (
	"io"

	"github.com/skycoin/skycoin-cx-forth/forth"
)

// Load translates p and compiles it into vm.
func Load(vm *forth.VM, p *Program) error {
	src, err := Translate(p)
	if err != nil {
		return err
	}
	return vm.Interpret(src)
}

// Run translates p, compiles it into a fresh VM and runs main.
func Run(p *Program, out io.Writer) error {
	vm := forth.New(out)
	if err := Load(vm, p); err != nil {
		return err
	}
	return vm.Interpret(WordName("main"))
}
