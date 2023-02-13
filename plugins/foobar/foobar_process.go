package main

import (
	"fmt"
	"github.com/a-peyrard/addon-manager/process"
)

type Process struct{}

func (p *Process) Run() (err error) {
	fmt.Printf("\n== foooooooobaaaaaaaar! ==\n\n")

	return
}

func NewProcess() process.Process {
	return &Process{}
}

func Foo() {
	fmt.Printf("\n== foooooooo! ==\n\n")
}
