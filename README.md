# Little Machines

These are toy implementations around small, esoteric languages, written in Go.

## Thue

A language built around string substitution.

https://en.wikipedia.org/wiki/Thue_(programming_language)

## Hello World

```
a::=hello
hello::=~Hello from Thue!
::=
a
```

### Interpreter

`github.com/jmhobbs/little-machines/thue` is an interpreter

```
usage: thue [options] <program file>

  -print-state
       print the state when execution stops
```

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
