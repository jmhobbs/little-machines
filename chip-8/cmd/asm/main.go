package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmhobbs/little-machines/chip-8/asm"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	ctr := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		originalLine := scanner.Text()
		ctr = ctr + 1

		line := removeComments(originalLine)
		if len(line) == 0 {
			continue
		}

		instruction, operands := parseLine(line)

		opcodes, err := asm.Encode(instruction, operands)
		if err != nil {
			log.Fatalf("%v on line %d", err, ctr)
		}

		_, err = out.Write(opcodes)
		if err != nil {
			log.Fatalf("error writing to outpur: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

/*
Strip comments out of a line.
*/
func removeComments(line string) string {
	segments := strings.SplitN(line, ";", 2)
	return strings.TrimSpace(segments[0])
}

/*
Split a line into instruction and operands.
Assumes comments are already removed.
Line is expected to match "<ins> <op>(, <op>)", ignoring extra whitespace.
*/
func parseLine(line string) (string, []string) {
	front, back, _ := strings.Cut(line, ",")
	instruction, operand, _ := strings.Cut(strings.TrimSpace(front), " ")

	operands := []string{}
	if len(operand) > 0 {
		operands = append(operands, strings.TrimSpace(strings.ToUpper(operand)))
	}

	backOperands := strings.Split(back, ",")
	for _, operand := range backOperands {
		operand = strings.TrimSpace(operand)
		if len(operand) > 0 {
			operands = append(operands, strings.ToUpper(operand))
		}
	}

	return strings.ToUpper(instruction), operands
}
