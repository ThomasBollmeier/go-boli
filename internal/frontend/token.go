package frontend

type TokenType int

const (
	TokLeftParen TokenType = iota
	TokRightParen
	TokLeftBrace
	TokRightBrace
	TokLeftBracket
	TokRightBracket
	TokQuotParen
	TokQuotBrace
	TokQuotBracket
	TokIdentifier
	TokVarParam
	TokSpread
	TokNil
	TokBoolean
	TokInteger
	TokReal
	TokRational
	TokString
	TokSymbol
	TokPlus
	TokMinus
	TokAsterisk
	TokSlash
	TokSlashSlash
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
	TokDefStruct
	TokSetBang
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
	'^': TokCaret,
	'%': TokPercentage,
}

var OpeningClosingPairs = map[TokenType]TokenType{
	TokLeftParen:   TokRightParen,
	TokLeftBracket: TokRightBracket,
	TokLeftBrace:   TokRightBrace,
	TokQuotParen:   TokRightParen,
	TokQuotBrace:   TokRightBrace,
	TokQuotBracket: TokRightBracket,
}

var Keywords = map[string]TokenType{
	"def":        TokDef,
	"def-struct": TokDefStruct,
	"if":         TokIf,
	"cond":       TokCond,
	"block":      TokBlock,
	"let":        TokLet,
	"#true":      TokBoolean,
	"#t":         TokBoolean,
	"#false":     TokBoolean,
	"#f":         TokBoolean,
	"nil":        TokNil,
	"and":        TokAnd,
	"or":         TokOr,
	"lambda":     TokLambda,
	"λ":          TokLambda,
	"set!":       TokSetBang,
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
