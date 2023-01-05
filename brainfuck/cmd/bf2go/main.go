package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var memorySize *uint = flag.Uint("memory-size", 300, "number of memory cells to use")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: bf2go [options]")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		flag.PrintDefaults()
	}

	flag.Parse()

	_, err := fmt.Fprintf(os.Stdout, `package main
	
import "fmt"
import "github.com/jmhobbs/little-machines/brainfuck/bf"

func read() byte {
	fmt.Print("input: ")
	var r rune
	_, err := fmt.Scanf("%%c", &r)
	if err != nil {
		panic(err)
	}
	return byte(r)
}

func main() {
	m := bf.NewMachine(%d)
`, *memorySize)

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		switch buf[0] {
		case '>':
			os.Stdout.Write([]byte("m.PtrIncr()\n"))
		case '<':
			os.Stdout.Write([]byte("m.PtrDecr()\n"))
		case '+':
			os.Stdout.Write([]byte("m.Incr()\n"))
		case '-':
			os.Stdout.Write([]byte("m.Decr()\n"))
		case '.':
			os.Stdout.Write([]byte("fmt.Print(string(m.Peek()))\n"))
		case ',':
			os.Stdout.Write([]byte("m.Put(read())\n"))
		case '[':
			os.Stdout.Write([]byte("for ; m.Peek() != 0; {\n"))
		case ']':
			os.Stdout.Write([]byte("}\n"))
		}
	}

	os.Stdout.Write([]byte("}\n"))
}
