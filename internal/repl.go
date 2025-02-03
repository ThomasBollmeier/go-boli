package internal

import (
	"fmt"
	"github.com/ergochat/readline"
	"go-boli/internal/frontend"
	ip "go-boli/internal/interpreter"
)

const Version = "0.4.3"

func Repl() {
	var input string
	var value ip.ValueObject

	fmt.Printf("(B)ollmeier's (O)wn (L)isp (I)mplementation - Version %s\n", Version)
	fmt.Println("Type ':q' to quit")
	fmt.Println()

	rl, err := readline.New("boλi> ")
	if err != nil {
		return
	}
	defer func(rl *readline.Instance) {
		_ = rl.Close()
	}(rl)

	interpreter := ip.NewInterpreter(nil)
	code := ""

	for {
		input, err = rl.Readline()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if input == ":q" {
			break
		}

		code += input
		if isToBeContinued(code) {
			continue
		}

		value, err = interpreter.Run(code)
		code = ""

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%v\n", value)
	}
}

func isToBeContinued(input string) bool {
	tokens := tokenize(input)
	level := 0
	for _, token := range tokens {
		switch token.Type {
		case frontend.TokLeftParen, frontend.TokLeftBracket, frontend.TokLeftBrace:
			level++
		case frontend.TokRightParen, frontend.TokRightBracket, frontend.TokRightBrace:
			level--
		default:
			continue
		}
	}
	return level > 0
}

func tokenize(code string) []*frontend.Token {
	var tokens []*frontend.Token
	scanner := frontend.NewScanner(frontend.NewCharStreamString(code))
	for {
		token, err := scanner.Advance()
		if err != nil {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens
}
