package frontend

type AstType int

const (
	AstInteger AstType = iota
	AstRational
	AstReal
	AstString
	AstVariable
	AstCall
	AstOperator
)

type AST struct {
	astType  AstType
	value    string
	children []*AST
	tokens   []*Token
	attrs    map[string]interface{}
}

func NewAST(astType AstType, value string) *AST {
	return &AST{
		astType:  astType,
		value:    value,
		children: make([]*AST, 0),
		tokens:   make([]*Token, 0),
		attrs:    make(map[string]interface{}),
	}
}

func NewASTAtom(astType AstType, token *Token) *AST {
	return &AST{
		astType: astType,
		value:   token.Lexeme,
		tokens:  []*Token{token},
		attrs:   make(map[string]interface{}),
	}

}

func NewASTNumber(token *Token) *AST {
	var astType AstType

	switch token.Type {
	case TokInteger:
		astType = AstInteger
	case TokRational:
		astType = AstRational
	case TokReal:
		astType = AstReal
	default:
		panic("unknown number type")
	}

	return NewASTAtom(astType, token)
}

func (ast *AST) AddChild(child *AST) {
	ast.children = append(ast.children, child)
	for _, token := range child.tokens {
		ast.tokens = append(ast.tokens, token)
	}
}

func (ast *AST) AddToken(token *Token) {
	ast.tokens = append(ast.tokens, token)
}
