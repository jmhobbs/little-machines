package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var printState *bool = flag.Bool("print-state", false, "print the state when execution stops")

	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: thue [options] <program file>")
		fmt.Fprintln(flag.CommandLine.Output(), "")
		flag.PrintDefaults()
	}

	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	m, err := New(f, os.Stdin)
	if err != nil {
		panic(err)
	}
	m.Run(*printState)
}

type machine struct {
	input *bufio.Scanner
	rules map[string]string
	state string
}

func New(program, input io.Reader) (*machine, error) {
	rules := map[string]string{}
	state := []string{}

	scanner := bufio.NewScanner(program)
	line := 0

	for scanner.Scan() {
		line++

		// allow empty lines
		if scanner.Text() == "" {
			continue
		}

		split := strings.SplitN(scanner.Text(), "::=", 2)

		if len(split) == 1 {
			return nil, fmt.Errorf("invalid program line %d, no ::=\n%q\n", line, scanner.Text())
		}

		// like the original interpreter, we treat "::=.*" as program end
		if split[0] == "" {
			break
		}
		rules[split[0]] = split[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for scanner.Scan() {
		state = append(state, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &machine{input: bufio.NewScanner(input), rules: rules, state: strings.Join(state, "")}, nil
}

func (m *machine) Run(printState bool) {
	var err error
	for {
		err = m.Step()
		if err != nil {
			if err == io.EOF {
				fmt.Println(m.state)
				os.Exit(0)
			}
			panic(err)
		}
	}
}

func (m *machine) Step() error {
	var (
		idx      int
		newState []string
	)

	for match, rule := range m.rules {
		idx = strings.Index(m.state, match)
		if idx != -1 {
			newState = []string{m.state[0:idx]}

			if rule[0] == '~' {
				// not in spec, but based on convention, if the output string is empty, print a newline
				if len(rule) == 1 {
					fmt.Println("")
				} else {
					fmt.Print(rule[1:])
				}
			} else if rule == ":" {
				// todo: read here
				if m.input.Scan() {
					newState = append(newState, m.input.Text())
				} else {
					return m.input.Err()
				}
			} else {
				newState = append(newState, rule)
			}
			newState = append(newState, m.state[idx+len(match):])

			m.state = strings.Join(newState, "")
			return nil
		}
	}
	return io.EOF
}
