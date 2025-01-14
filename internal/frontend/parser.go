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
		var child *AST

		token, err := p.scanner.Advance()
		if err != nil {
			return nil, err
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
				child, err = p.parseDefinition(token)
			default:
				child, err = p.parseExpr(token)
			}
		default:
			child, err = p.parseExpr(token)
		}

		if err != nil {
			return nil, err
		}

		program.AddChild(child)
	}

	return program, nil
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
	_, err := p.expect(TokDef) // scan def keyword
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
	ret.AddToken(identifier)
	ret.AddChild(value)
	ret.AddToken(end)

	return ret, nil
}

func (p *Parser) parseIfExpr(start *Token) (*AST, error) {

	closingType := OpeningClosingPairs[start.Type]
	_, err := p.expect(TokIf)
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
	ret.AddChild(condition)
	ret.AddChild(consequent)
	ret.AddChild(alternative)
	ret.AddToken(end)

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
