package internal

type TokenType int

const (
	TokLeftParen TokenType = iota
	TokRightParen
	TokLeftBrace
	TokRightBrace
	TokLeftBracket
	TokRightBracket
	TokIdentifier
	TokInteger
	TokReal
	TokRational
	TokString
	TokPlus
	TokMinus
	TokAsterisk
	TokSlash
)

type Token struct {
	Type   TokenType
	Lexeme string
	Row    int
	Col    int
}
