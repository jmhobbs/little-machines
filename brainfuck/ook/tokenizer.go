package ook

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jmhobbs/little-machines/brainfuck/bf"
)

var tokenMap map[string]bf.Token = map[string]bf.Token{
	"Ook. Ook?": bf.MoveRight,
	"Ook? Ook.": bf.MoveLeft,
	"Ook. Ook.": bf.Increment,
	"Ook! Ook!": bf.Decrement,
	"Ook! Ook.": bf.Output,
	"Ook. Ook!": bf.Input,
	"Ook! Ook?": bf.JumpForward,
	"Ook? Ook!": bf.JumpBackward,
	"Ook? Ook?": bf.Noop,
}

func Tokenize(in io.Reader) ([]bf.Token, error) {
	tokens := []bf.Token{}

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)

	var (
		prefix, token string
		bfToken       bf.Token
		ok            bool
	)

	for scanner.Scan() {
		if prefix == "" {
			prefix = scanner.Text()
		} else {
			token = prefix + " " + scanner.Text()
			prefix = ""

			bfToken, ok = tokenMap[token]
			if !ok {
				return nil, fmt.Errorf("invalid command %q", token)
			}
			tokens = append(tokens, bfToken)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	return tokens, nil

}
