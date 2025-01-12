package frontend

import (
	"errors"
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
	s.skipWhitespace()

	row := s.row
	col := s.col

	ch, err := s.advanceChar()
	if err != nil {
		return nil, err
	}

	tokType, ok := SingleCharTokens[ch]
	if ok {
		return NewToken(tokType, string(ch), row, col), nil
	}

	if unicode.IsDigit(ch) {
		return s.scanNumber(ch, row, col)
	}

	if ch == '"' {
		return s.scanString(ch, row, col)
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
	if !ok {
		tokenType = TokIdentifier
	}

	return NewToken(tokenType, lexeme, row, col), nil
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
	invalidChars := []rune{'(', ')', '{', '}', '[', ']', ':'}
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
