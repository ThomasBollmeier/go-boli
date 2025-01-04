package frontend

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
	TokCaret
	TokPercentage
)

var SingleCharTokens = map[rune]TokenType{
	'(': TokLeftParen,
	')': TokRightParen,
	'{': TokLeftBrace,
	'}': TokRightBrace,
	'[': TokLeftBracket,
	']': TokRightBracket,
	'+': TokPlus,
	'-': TokMinus,
	'*': TokAsterisk,
	'/': TokSlash,
	'^': TokCaret,
	'%': TokPercentage,
}

type Token struct {
	Type   TokenType
	Lexeme string
	Row    int
	Col    int
}

func NewToken(typ TokenType, lexeme string, row, col int) *Token {
	return &Token{
		Type:   typ,
		Lexeme: lexeme,
		Row:    row,
		Col:    col,
	}
}
