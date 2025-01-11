package frontend

import "errors"

type Parser struct {
	scanner *Scanner
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(stream CharStream) (*AST, error) {
	p.scanner = NewScanner(stream)
	token, err := p.scanner.Advance()
	if err != nil {
		return nil, err
	}

	return p.parseExpr(token)
}

func (p *Parser) parseExpr(token *Token) (*AST, error) {
	var err error

	if token == nil {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}
	}

	switch token.Type {
	case TokInteger, TokRational, TokReal:
		return NewASTNumber(token), nil
	case TokString:
		return NewASTAtom(AstString, token), nil
	case TokIdentifier:
		return NewASTAtom(AstVariable, token), nil
	case TokLeftParen, TokLeftBracket, TokLeftBrace:
		return p.parseCall(token)
	case TokPlus, TokMinus, TokAsterisk, TokSlash, TokCaret, TokPercentage:
		return NewASTAtom(AstOperator, token), nil
	default:
		return nil, errors.New("unknown expression")
	}
}

func (p *Parser) parseCall(start *Token) (*AST, error) {
	var calleeAst, argAst *AST
	var token *Token
	var err error

	callee, err := p.scanner.Advance()
	if err != nil {
		return nil, err
	}
	switch callee.Type {
	case TokPlus, TokMinus, TokAsterisk, TokSlash,
		TokCaret, TokPercentage:
		calleeAst = NewASTAtom(AstOperator, callee)
	default:
		return nil, errors.New("unknown callee type")
	}

	callAst := NewAST(AstCall, "")
	callAst.AddToken(start)
	callAst.AddChild(calleeAst)

	closingType := OpeningClosingPairs[start.Type]

	for {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}
		if token.Type == closingType {
			callAst.AddToken(token)
			break
		}
		argAst, err = p.parseExpr(token)
		if err != nil {
			return nil, err
		}
		callAst.AddChild(argAst)
	}

	return callAst, nil
}
