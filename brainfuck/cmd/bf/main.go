package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/jmhobbs/little-machines/brainfuck/bf"
)

func main() {
	var memorySize *uint = flag.Uint("memory-size", 300, "number of memory cells to use")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: bf [options] <program file>")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		flag.PrintDefaults()
	}

	flag.Parse()

	program, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	i := bf.NewInterpreter(*memorySize, program)
	i.Run()
}
