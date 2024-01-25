package main

import (
	"io/ioutil"
	"os"

	"github.com/jmhobbs/little-machines/chip-8/chip8"
)

func main() {
	pgm, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	m, err := chip8.New(pgm)
	if err != nil {
		panic(err)
	}

	if err := m.Run(); err != nil {
		panic(err)
	}
}
