package frontend

import (
	"errors"
	"fmt"
)

type Parser struct {
	scanner *BufferedStream[*Token]
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(stream CharStream) (*AST, error) {
	p.scanner = NewBufferedStream(NewScanner(stream))
	program := NewAST(AstProgram, "")

	for p.scanner.HasNext() {
		child, err := p.parseDefOrExpr(nil)
		if err != nil {
			return nil, err
		}
		program.AddChild(child)
	}

	return program, nil
}

func (p *Parser) parseDefOrExpr(token *Token) (*AST, error) {
	var err error

	if token == nil {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}
	}

	var nextToken *Token
	switch token.Type {
	case TokLeftParen, TokLeftBrace, TokLeftBracket:
		nextToken, err = p.scanner.Peek()
		if err != nil {
			return nil, err
		}
		switch nextToken.Type {
		case TokDef:
			return p.parseDefinition(token)
		default:
			return p.parseExpr(token)
		}
	default:
		return p.parseExpr(token)
	}
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
	case TokNil:
		return NewASTAtom(AstNil, token), nil
	case TokBoolean:
		return NewASTAtom(AstBoolean, token), nil
	case TokInteger, TokRational, TokReal:
		return NewASTNumber(token), nil
	case TokString:
		return NewASTAtom(AstString, token), nil
	case TokIdentifier:
		return NewASTAtom(AstVariable, token), nil
	case TokLeftParen, TokLeftBracket, TokLeftBrace:
		nextToken, errPeek := p.scanner.Peek()
		if errPeek != nil {
			return nil, errPeek
		}
		switch nextToken.Type {
		case TokIf:
			return p.parseIfExpr(token)
		case TokCond:
			return p.parseCondExpr(token)
		case TokBlock:
			return p.parseBlock(token)
		case TokLet:
			return p.parseLetExpr(token)
		case TokLambda:
			return p.parseLambda(token)
		case TokAnd, TokOr:
			return p.parseLogicalOp(token)
		default:
			return p.parseCall(token)
		}
	case TokPlus, TokMinus, TokAsterisk, TokSlash, TokCaret, TokPercentage:
		return NewASTAtom(AstOperator, token), nil
	case TokEqual, TokGreater, TokGreaterEq, TokLess, TokLessEq:
		return NewASTAtom(AstComparisonOp, token), nil
	default:
		return nil, errors.New("unknown expression")
	}
}

func (p *Parser) parseLogicalOp(start *Token) (*AST, error) {
	var token *Token
	var expr *AST
	var err error

	token, err = p.expect(TokAnd, TokOr)
	if err != nil {
		return nil, err
	}

	var astType AstType
	if token.Type == TokAnd {
		astType = AstConjunction
	} else {
		astType = AstDisjunction
	}

	ret := NewAST(astType, "")
	ret.AddToken(start)
	ret.AddToken(token)

	closingType := OpeningClosingPairs[start.Type]

	for {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}
		if token.Type == closingType {
			ret.AddToken(token)
			break
		}
		expr, err = p.parseExpr(token)
		if err != nil {
			return nil, err
		}
		ret.AddChild(expr)
	}

	return ret, nil
}

func (p *Parser) parseCall(start *Token) (*AST, error) {
	var calleeAst, argAst *AST
	var token *Token
	var err error

	callee, err := p.scanner.Advance()
	if err != nil {
		return nil, err
	}
	calleeAst, err = p.parseExpr(callee)

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

func (p *Parser) parseDefinition(start *Token) (*AST, error) {
	closingType := OpeningClosingPairs[start.Type]
	tokenDef, err := p.expect(TokDef) // scan def keyword
	if err != nil {
		return nil, err
	}

	var identifier *Token
	identifier, err = p.expect(TokIdentifier)
	if err != nil {
		return nil, err
	}

	var value *AST
	value, err = p.parseExpr(nil)
	if err != nil {
		return nil, err
	}

	var end *Token
	end, err = p.expect(closingType)
	if err != nil {
		return nil, err
	}

	ret := NewAST(AstDefinition, identifier.Lexeme)
	ret.AddToken(start)
	ret.AddToken(tokenDef)
	ret.AddToken(identifier)
	ret.AddChild(value)
	ret.AddToken(end)

	return ret, nil
}

func (p *Parser) parseIfExpr(start *Token) (*AST, error) {

	closingType := OpeningClosingPairs[start.Type]
	tokenIf, err := p.expect(TokIf)
	if err != nil {
		return nil, err
	}

	var condition *AST
	condition, err = p.parseExpr(nil)
	if err != nil {
		return nil, err
	}

	var consequent *AST
	consequent, err = p.parseExpr(nil)
	if err != nil {
		return nil, err
	}

	var alternative *AST
	alternative, err = p.parseExpr(nil)
	if err != nil {
		return nil, err
	}

	var end *Token
	end, err = p.expect(closingType)
	if err != nil {
		return nil, err
	}

	ret := NewAST(AstIfExpression, "")
	ret.AddToken(start)
	ret.AddToken(tokenIf)
	ret.AddChild(condition)
	ret.AddChild(consequent)
	ret.AddChild(alternative)
	ret.AddToken(end)

	return ret, nil
}

func (p *Parser) parseCondExpr(start *Token) (*AST, error) {
	var token *Token
	var err error
	var prevIf *AST = nil
	var clauseCond *AST
	var clauseExpr *AST

	ret := NewAST(AstIfExpression, "")
	ret.AddToken(start)

	closingType := OpeningClosingPairs[start.Type]
	token, err = p.expect(TokCond)
	if err != nil {
		return nil, err
	}

	ret.AddToken(token)

	currIf := ret

	for {
		token, err = p.expect(closingType, TokLeftParen, TokLeftBrace, TokLeftBracket)
		if err != nil {
			return nil, err
		}

		if token.Type == closingType {
			ret.AddToken(token)
			break
		}

		currIf.AddToken(token)
		clauseClosingType := OpeningClosingPairs[token.Type]

		clauseCond, err = p.parseExpr(nil)
		if err != nil {
			return nil, err
		}
		currIf.AddChild(clauseCond)

		clauseExpr, err = p.parseExpr(nil)
		if err != nil {
			return nil, err
		}
		currIf.AddChild(clauseExpr)

		token, err = p.expect(clauseClosingType)
		if err != nil {
			return nil, err
		}
		currIf.AddToken(token)

		nextIf := NewAST(AstIfExpression, "")
		currIf.AddChild(nextIf)
		prevIf = currIf
		currIf = nextIf
	}

	if len(ret.GetChildren()) == 0 {
		return nil, errors.New("cond expression has no clauses")
	}

	if prevIf != nil {
		prevIf.ReplaceLastChild(NewAST(AstNil, ""))
	}

	return ret, nil
}

func (p *Parser) parseBlock(token *Token) (*AST, error) {
	var child *AST

	ret := NewAST(AstBlock, "")
	ret.AddToken(token)

	tokenBlock, err := p.expect(TokBlock)
	if err != nil {
		return nil, err
	}
	ret.AddToken(tokenBlock)

	closingType := OpeningClosingPairs[token.Type]

	for {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}

		if token.Type == closingType {
			ret.AddToken(token)
			break
		}

		child, err = p.parseDefOrExpr(token)
		if err != nil {
			return nil, err
		}

		ret.AddChild(child)
	}

	return ret, nil
}

func (p *Parser) parseLetExpr(token *Token) (*AST, error) {
	var definitions []*AST
	var tokenDefStart, tokenDefEnd *Token
	var child *AST

	ret := NewAST(AstBlock, "")
	ret.AddToken(token)

	closingType := OpeningClosingPairs[token.Type]

	token, err := p.expect(TokLet)
	if err != nil {
		return nil, err
	}
	ret.AddToken(token)

	tokenDefStart, definitions, tokenDefEnd, err = p.parseLetDefinitions()
	if err != nil {
		return nil, err
	}
	ret.AddToken(tokenDefStart)
	for _, definition := range definitions {
		ret.AddChild(definition)
	}
	ret.AddToken(tokenDefEnd)

	for {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}

		if token.Type == closingType {
			ret.AddToken(token)
			break
		}

		child, err = p.parseExpr(token)
		if err != nil {
			return nil, err
		}

		ret.AddChild(child)
	}

	return ret, nil
}

func (p *Parser) parseLetDefinitions() (*Token, []*AST, *Token, error) {
	var defs []*AST
	var def *AST
	var start, end, token *Token
	var err error

	start, err = p.expect(TokLeftParen, TokLeftBrace, TokLeftBracket)
	if err != nil {
		return nil, nil, nil, err
	}
	endType := OpeningClosingPairs[start.Type]

	for {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, nil, nil, err
		}
		if token.Type == endType {
			end = token
			break
		}
		def, err = p.parseLetDefinition(token)
		if err != nil {
			return nil, nil, nil, err
		}
		defs = append(defs, def)
	}

	return start, defs, end, err
}

func (p *Parser) parseLetDefinition(start *Token) (*AST, error) {

	endType := OpeningClosingPairs[start.Type]

	identifier, err := p.expect(TokIdentifier)
	if err != nil {
		return nil, err
	}

	value, err := p.parseExpr(nil)
	if err != nil {
		return nil, err
	}

	end, err := p.expect(endType)
	if err != nil {
		return nil, err
	}

	ret := NewAST(AstDefinition, identifier.Lexeme)
	ret.AddToken(start)
	ret.AddToken(identifier)
	ret.AddChild(value)
	ret.AddToken(end)

	return ret, nil
}

func (p *Parser) parseLambda(start *Token) (*AST, error) {
	ret := NewAST(AstLambda, "")
	ret.AddToken(start)

	endType := OpeningClosingPairs[start.Type]

	token, err := p.expect(TokLambda)
	if err != nil {
		return nil, err
	}
	ret.AddToken(token)

	token, err = p.expect(TokLeftParen, TokLeftBrace, TokLeftBracket)
	if err != nil {
		return nil, err
	}
	ret.AddToken(token)

	params := NewAST(AstParameters, "")
	endParamsType := OpeningClosingPairs[token.Type]

	var tokenEndParams *Token
paramsLoop:
	for {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}
		switch token.Type {
		case TokIdentifier:
			params.AddChild(NewASTAtom(AstVariable, token))
		case endParamsType:
			tokenEndParams = token
			break paramsLoop
		default:
			return nil, errors.New("expected identifier as parameter")
		}
	}

	ret.AddChild(params)
	ret.AddToken(tokenEndParams)

	block := NewAST(AstBlock, "")
	var tokenEnd *Token
	var child *AST

	for {
		token, err = p.scanner.Advance()
		if err != nil {
			return nil, err
		}
		if token.Type == endType {
			tokenEnd = token
			break
		}

		child, err = p.parseDefOrExpr(token)
		if err != nil {
			return nil, err
		}
		block.AddChild(child)
	}

	ret.AddChild(block)
	ret.AddToken(tokenEnd)

	return ret, nil
}

func (p *Parser) expect(tokenTypes ...TokenType) (*Token, error) {
	ret, err := p.scanner.Advance()
	if err != nil {
		return nil, err
	}
	for _, tokenType := range tokenTypes {
		if tokenType == ret.Type {
			return ret, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("expected %v but got %v", tokenTypes, ret.Type))
}
