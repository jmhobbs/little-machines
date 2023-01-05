# Little Machines

These are toy implementations around small, esoteric languages, written in Go.

## Brainfuck

A small language with only 8 commands, encoded as the characters `>`, `<`, `+`, `-`, `.`, `,`, `[`, and `]`

https://en.wikipedia.org/wiki/Brainfuck

### Interpreter

`github.com/jmhobbs/little-machines/brainfuck/cmd/bf` is an interpreter

```
usage: bf [options] <program file>

  -memory-size uint
    	number of memory cells to use (default 300)
```

### Go Encoder

`github.com/jmhobbs/little-machines/brainfuck/cmd/bf2go` takes a program on stdin and converts it to Go source on stdout

```
usage: bf2go [options]

  -memory-size uint
    	number of memory cells to use (default 300)
```
