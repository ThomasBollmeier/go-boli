package frontend

import (
	"errors"
	"fmt"
	"unicode"
)

type Scanner struct {
	inStream BufferedStream[rune]
	row      int
	col      int
}

func NewScanner(stream CharStream) *Scanner {
	return &Scanner{
		inStream: *NewBufferedStream(stream),
		row:      1,
		col:      1,
	}
}

func (s *Scanner) Advance() (*Token, error) {
	var row, col int
	var ch rune
	var err error

	for {
		s.skipWhitespace()

		row = s.row
		col = s.col

		ch, err = s.advanceChar()
		if err != nil {
			return nil, err
		}

		if ch == ';' { // start of a line comment
			err = s.skipLineComment()
			if err != nil {
				return nil, err
			}
			continue
		}

		if ch == '#' {
			var nextChar rune
			nextChar, err = s.inStream.Peek()
			if err == nil {
				if nextChar == '|' {
					_, _ = s.Advance()
					err = s.skipBlockComment()
					if err != nil {
						return nil, err
					}
					continue
				}
				if nextChar == '!' {
					err = s.skipLineComment()
					if err != nil {
						return nil, err
					}
					continue
				}
			}
		}
		break
	}

	tokType, ok := SingleCharTokens[ch]
	if ok {
		return NewToken(tokType, string(ch), row, col), nil
	}

	switch ch {
	case '+', '-':
		nextCh, errNext := s.inStream.Peek()
		if ch == '+' {
			tokType = TokPlus
		} else {
			tokType = TokMinus
		}
		if errNext != nil || !unicode.IsDigit(nextCh) {
			return NewToken(tokType, string(ch), row, col), nil
		} else {
			return s.scanNumber(ch, row, col)
		}
	case '=':
		return NewToken(TokEqual, "=", row, col), nil
	case '>':
		nextCh, errNext := s.inStream.Peek()
		if errNext == nil && nextCh == '=' {
			_, _ = s.advanceChar()
			return NewToken(TokGreaterEq, string(ch)+string(nextCh), row, col), nil
		} else {
			return NewToken(TokGreater, string(ch), row, col), nil
		}
	case '<':
		nextCh, errNext := s.inStream.Peek()
		if errNext == nil && nextCh == '=' {
			_, _ = s.advanceChar()
			return NewToken(TokLessEq, string(ch)+string(nextCh), row, col), nil
		} else {
			return NewToken(TokLess, string(ch), row, col), nil
		}
	case '.':
		nextChars := s.inStream.PeekMany(3)
		if len(nextChars) == 3 &&
			nextChars[0] == '.' &&
			nextChars[1] == '.' &&
			s.isValidFirstIdentChar(nextChars[2]) {

			for i := 0; i < 3; i++ {
				_, _ = s.advanceChar()
			}

			var ident *Token
			ident, err = s.scanIdentifier(nextChars[2], row, col)
			if err != nil {
				return nil, err
			}
			if ident.Type != TokIdentifier {
				return nil, errors.New("spread expects identifier")
			}
			return NewToken(TokSpread, "..."+ident.Lexeme, row, col), nil
		} else {
			return NewToken(TokDot, ".", row, col), nil
		}
	case '"':
		return s.scanString(ch, row, col)
	case '\'':
		var nextChar rune
		nextChar, err = s.inStream.Peek()
		if err != nil {
			return nil, err
		}
		switch nextChar {
		case '(':
			_, _ = s.advanceChar()
			return NewToken(TokQuotParen, string(ch)+string(nextChar), row, col), nil
		case '{':
			_, _ = s.advanceChar()
			return NewToken(TokQuotBrace, string(ch)+string(nextChar), row, col), nil
		case '[':
			_, _ = s.advanceChar()
			return NewToken(TokQuotBracket, string(ch)+string(nextChar), row, col), nil
		default:
			return s.scanSymbol(ch, row, col)
		}
	}

	if unicode.IsDigit(ch) {
		return s.scanNumber(ch, row, col)
	}

	if s.isValidFirstIdentChar(ch) {
		return s.scanIdentifier(ch, row, col)
	}

	return nil, errors.New("no token could be found")
}

func (s *Scanner) HasNext() bool {
	s.skipWhitespace()
	return s.inStream.HasNext()
}

func (s *Scanner) scanSymbol(ch rune, row, col int) (*Token, error) {
	lexeme := string(ch)
	first := true
	for {
		nextCh, err := s.inStream.Peek()
		if err != nil {
			if first {
				return nil, err
			} else {
				break
			}
		}
		if first && !s.isValidFirstIdentChar(nextCh) {
			return nil, fmt.Errorf("invalid start of symbol: %c", nextCh)
		}
		if !first && !s.isValidIdentChar(nextCh) {
			break
		}
		first = false
		_, _ = s.advanceChar()
		lexeme += string(nextCh)
	}

	return NewToken(TokSymbol, lexeme, row, col), nil
}

func (s *Scanner) scanIdentifier(ch rune, row, col int) (*Token, error) {
	lexeme := string(ch)
	for {
		nextCh, err := s.inStream.Peek()
		if err != nil {
			break
		}
		if !s.isValidIdentChar(nextCh) {
			break
		}
		_, _ = s.advanceChar()
		lexeme += string(nextCh)
	}

	tokenType, ok := Keywords[lexeme]
	if ok {
		return NewToken(tokenType, lexeme, row, col), nil
	}

	nextChars := s.inStream.PeekMany(3)
	if string(nextChars) == "..." {
		for i := 0; i < 3; i++ {
			_, _ = s.advanceChar()
		}
		lexeme += "..."
		return NewToken(TokVarParam, lexeme, row, col), nil
	}

	return NewToken(TokIdentifier, lexeme, row, col), nil
}

func (s *Scanner) isValidFirstIdentChar(ch rune) bool {
	if unicode.IsLetter(ch) {
		return true
	}
	if ch == '#' {
		return true
	}
	return false
}

func (s *Scanner) isValidIdentChar(ch rune) bool {
	if unicode.IsSpace(ch) {
		return false
	}
	invalidChars := []rune{'(', ')', '{', '}', '[', ']', '.', ';'}
	for _, invalidChar := range invalidChars {
		if ch == invalidChar {
			return false
		}
	}

	return true
}

func (s *Scanner) scanString(firstChar rune, row, col int) (*Token, error) {
	lexeme := string(firstChar)
	prevChar := rune(0)

	for {
		ch, err := s.advanceChar()
		if err != nil {
			return nil, err
		}
		lexeme += string(ch)
		if ch == firstChar && prevChar != '\\' {
			break
		}
		prevChar = ch
	}

	return NewToken(TokString, lexeme, row, col), nil
}

func (s *Scanner) scanNumber(firstDigit rune, row, col int) (*Token, error) {
	tokType := TokInteger
	lexeme := string(firstDigit)

	for {
		ch, err := s.inStream.Peek()
		if err != nil {
			break
		}
		if unicode.IsDigit(ch) {
			lexeme += string(ch)
			_, _ = s.advanceChar()
			continue
		}
		switch ch {
		case ',', '/':
			if ch == ',' {
				tokType = TokReal
			} else {
				tokType = TokRational
			}
			lexeme += string(ch)
			_, _ = s.advanceChar()
		}
		break
	}

	if tokType == TokInteger {
		return NewToken(tokType, lexeme, row, col), nil
	}

	nextDigits := ""
	for {
		ch, err := s.inStream.Peek()
		if err != nil {
			break
		}
		if unicode.IsDigit(ch) {
			nextDigits += string(ch)
			_, _ = s.advanceChar()
			continue
		}
		break
	}

	if len(nextDigits) == 0 {
		return nil, errors.New("invalid number")
	}

	lexeme += nextDigits

	return NewToken(tokType, lexeme, row, col), nil
}

func (s *Scanner) advanceChar() (rune, error) {
	ch, err := s.inStream.Advance()
	if err != nil {
		return 0, err
	}
	if ch != '\n' {
		s.col++
	} else {
		s.row++
		s.col = 1
	}
	return ch, nil
}

func (s *Scanner) skipWhitespace() {
	for {
		ch, err := s.inStream.Peek()
		if err != nil {
			break
		}
		if !unicode.IsSpace(ch) {
			break
		}
		_, _ = s.advanceChar()
	}
}

func (s *Scanner) skipLineComment() error {
	for {
		ch, err := s.advanceChar()
		if err != nil {
			return err
		}
		if string(ch) == "\n" {
			return nil
		}
	}
}

func (s *Scanner) skipBlockComment() error {
	for {
		nextChars := string(s.inStream.PeekMany(2))
		if nextChars == "|#" {
			_, _ = s.advanceChar()
			_, _ = s.advanceChar()
			return nil
		}
		_, err := s.advanceChar()
		if err != nil {
			return err
		}
	}
}
