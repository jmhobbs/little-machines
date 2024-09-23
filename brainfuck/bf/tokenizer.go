package bf

import (
	"errors"
	"io"
)

type Token byte

const (
	MoveRight    Token = 62 // >
	MoveLeft     Token = 60 // <
	Increment    Token = 43 // +
	Decrement    Token = 45 // -
	Output       Token = 46 // .
	Input        Token = 44 // ,
	JumpForward  Token = 91 // [
	JumpBackward Token = 93 // ]
)

type Tokenizer func(io.Reader) ([]Token, error)

func Tokenize(in io.Reader) ([]Token, error) {
	b := make([]byte, 1)
	tokens := []Token{}

	for {
		_, err := in.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return tokens, err
		}

		switch Token(b[0]) {
		case MoveRight:
			fallthrough
		case MoveLeft:
			fallthrough
		case Increment:
			fallthrough
		case Decrement:
			fallthrough
		case Output:
			fallthrough
		case Input:
			fallthrough
		case JumpForward:
			fallthrough
		case JumpBackward:
			tokens = append(tokens, Token(b[0]))
		default:
			// noop
		}
	}

	return tokens, nil
}
