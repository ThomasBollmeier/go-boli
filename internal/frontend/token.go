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
	TokDot
	TokEqual
	TokGreater
	TokGreaterEq
	TokLess
	TokLessEq
	TokAnd
	TokOr
	TokDef
	TokIf
	TokCond
	TokBlock
	TokLet
	TokLambda
)

var SingleCharTokens = map[rune]TokenType{
	'(': TokLeftParen,
	')': TokRightParen,
	'{': TokLeftBrace,
	'}': TokRightBrace,
	'[': TokLeftBracket,
	']': TokRightBracket,
	'*': TokAsterisk,
	'/': TokSlash,
	'^': TokCaret,
	'%': TokPercentage,
	'.': TokDot,
	'λ': TokLambda,
}

var OpeningClosingPairs = map[TokenType]TokenType{
	TokLeftParen:   TokRightParen,
	TokLeftBracket: TokRightBracket,
	TokLeftBrace:   TokRightBrace,
}

var Keywords = map[string]TokenType{
	"def":    TokDef,
	"if":     TokIf,
	"cond":   TokCond,
	"block":  TokBlock,
	"let":    TokLet,
	"#true":  TokBoolean,
	"#t":     TokBoolean,
	"#false": TokBoolean,
	"#f":     TokBoolean,
	"nil":    TokNil,
	"and":    TokAnd,
	"or":     TokOr,
	"lambda": TokLambda,
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
