package frontend

import "unicode"

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

func (s *Scanner) Advance() (Token, error) {
	s.skipWhitespace()

	row := s.row
	col := s.col

	ch, err := s.inStream.Advance()
	if err != nil {
		return *NewToken(TokEndOfStream, "", row, col), err
	}

	tokType, ok := SingleCharTokens[ch]
	if ok {
		return *NewToken(tokType, string(ch), row, col), nil
	}

	panic("unimplemented")
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
