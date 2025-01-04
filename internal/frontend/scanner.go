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

	ch, err := s.inStream.Advance()
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

	return nil, errors.New("no token could be found")
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
			_, _ = s.advance()
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
			_, _ = s.advance()
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
			_, _ = s.advance()
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

func (s *Scanner) advance() (rune, error) {
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
		_, _ = s.advance()
	}
}
