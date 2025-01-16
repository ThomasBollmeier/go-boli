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
	TokNil
	TokBoolean
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
	TokEqual
	TokGreater
	TokGreaterEq
	TokLess
	TokLessEq
	TokAnd
	TokOr
	TokDef
	TokIf
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

var OpeningClosingPairs = map[TokenType]TokenType{
	TokLeftParen:   TokRightParen,
	TokLeftBracket: TokRightBracket,
	TokLeftBrace:   TokRightBrace,
}

var Keywords = map[string]TokenType{
	"def":    TokDef,
	"if":     TokIf,
	"#true":  TokBoolean,
	"#t":     TokBoolean,
	"#false": TokBoolean,
	"#f":     TokBoolean,
	"nil":    TokNil,
	"and":    TokAnd,
	"or":     TokOr,
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
