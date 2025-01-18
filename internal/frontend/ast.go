package frontend

type AstType string

const (
	AstNil          AstType = "Nil"
	AstBoolean      AstType = "Boolean"
	AstInteger      AstType = "Integer"
	AstRational     AstType = "Rational"
	AstReal         AstType = "Real"
	AstString       AstType = "String"
	AstVariable     AstType = "Variable"
	AstCall         AstType = "Call"
	AstOperator     AstType = "Operator"
	AstComparisonOp AstType = "ComparisonOp"
	AstDisjunction  AstType = "Disjunction"
	AstConjunction  AstType = "Conjunction"
	AstDefinition   AstType = "Definition"
	AstIfExpression AstType = "IfExpression"
	AstProgram      AstType = "Program"
)

type AST struct {
	astType          AstType
	value            string
	tokensOrChildren []TokenOrChild
	attrs            map[string]interface{}
}

type TokenOrChild struct {
	token *Token
	child *AST
}

func NewAST(astType AstType, value string) *AST {
	return &AST{
		astType:          astType,
		value:            value,
		tokensOrChildren: make([]TokenOrChild, 0),
		attrs:            make(map[string]interface{}),
	}
}

func NewASTAtom(astType AstType, token *Token) *AST {
	return &AST{
		astType:          astType,
		value:            token.Lexeme,
		tokensOrChildren: []TokenOrChild{{token, nil}},
		attrs:            make(map[string]interface{}),
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
	ast.tokensOrChildren = append(ast.tokensOrChildren, TokenOrChild{child: child})
}

func (ast *AST) ReplaceLastChild(newChild *AST) {
	lastIdx := -1
	for idx, tokenOrChild := range ast.tokensOrChildren {
		if tokenOrChild.child != nil {
			lastIdx = idx
		}
	}
	if lastIdx != -1 {
		ast.tokensOrChildren[lastIdx] = TokenOrChild{child: newChild}
	}
}

func (ast *AST) AddToken(token *Token) {
	ast.tokensOrChildren = append(ast.tokensOrChildren, TokenOrChild{token: token})
}

func (ast *AST) GetType() AstType {
	return ast.astType
}

func (ast *AST) GetValue() string {
	return ast.value
}

func (ast *AST) GetChildren() []*AST {
	var ret []*AST
	for _, tokenOrChild := range ast.tokensOrChildren {
		if tokenOrChild.child != nil {
			ret = append(ret, tokenOrChild.child)
		}
	}
	return ret
}

func (ast *AST) GetAttributes() map[string]interface{} {
	return ast.attrs
}

func (ast *AST) GetLexemes() string {
	ret := ""
	for _, tokenOrChild := range ast.tokensOrChildren {
		if len(ret) > 0 {
			ret += " "
		}
		if tokenOrChild.child != nil {
			ret += tokenOrChild.child.GetLexemes()
		} else {
			ret += tokenOrChild.token.Lexeme
		}
	}
	return ret
}
