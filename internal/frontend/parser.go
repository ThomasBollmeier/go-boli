package frontend

import (
	"errors"
	"fmt"
)

type Parser struct {
	scanner         *BufferedStream[*Token]
	callArgsNesting int
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

	findTailCalls(program)

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
		case TokDefStruct:
			return p.parseStructDefinition(token)
		case TokSetBang:
			return p.parseVarChange(token)
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
	case TokSymbol:
		return NewASTAtom(AstSymbol, token), nil
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
	case TokQuotParen, TokQuotBrace, TokQuotBracket:
		return p.parseQuote(token)
	case TokPlus, TokMinus, TokAsterisk, TokSlash,
		TokSlashSlash, TokCaret, TokPercentage:
		return NewASTAtom(AstOperator, token), nil
	case TokEqual, TokGreater, TokGreaterEq, TokLess, TokLessEq:
		return NewASTAtom(AstComparisonOp, token), nil
	case TokSpread:
		if p.callArgsNesting > 0 {
			return NewASTAtom(AstSpread, token), nil
		} else {
			return nil, errors.New("spreads are only allowed as call args")
		}
	default:
		return nil, errors.New("unknown expression: '" + token.Lexeme + "'")
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
	var first, argAst *AST
	var token *Token
	var err error

	firstToken, err := p.scanner.Advance()
	if err != nil {
		return nil, err
	}
	first, err = p.parseExpr(firstToken)
	if err != nil {
		return nil, err
	}

	token, err = p.scanner.Peek()
	if err != nil {
		return nil, err
	}

	if token.Type == TokDot {
		return p.parsePair(start, first)
	}

	closingType := OpeningClosingPairs[start.Type]

	callAst := NewAST(AstCall, "")
	callAst.AddToken(start)
	callAst.AddChild(first)

	p.callArgsNesting++

	for {
		token, err = p.scanner.Advance()
		if err != nil {
			p.callArgsNesting--
			return nil, err
		}
		if token.Type == closingType {
			callAst.AddToken(token)
			break
		}
		argAst, err = p.parseExpr(token)
		if err != nil {
			p.callArgsNesting--
			return nil, err
		}
		callAst.AddChild(argAst)
	}

	p.callArgsNesting--

	return callAst, nil
}

func (p *Parser) parsePair(start *Token, first *AST) (*AST, error) {
	var err error
	var dot, end *Token
	var second *AST

	endType := OpeningClosingPairs[start.Type]

	dot, err = p.expect(TokDot)
	if err != nil {
		return nil, err
	}

	second, err = p.parseExpr(nil)
	if err != nil {
		return nil, err
	}

	end, err = p.expect(endType)
	if err != nil {
		return nil, err
	}

	ret := NewAST(AstPair, "")
	ret.AddToken(start)
	ret.AddChild(first)
	ret.AddToken(dot)
	ret.AddChild(second)
	ret.AddToken(end)

	return ret, nil
}

func (p *Parser) parseQuote(start *Token) (*AST, error) {
	ret := NewAST(AstCall, "")
	ret.AddToken(start)
	ret.AddChild(NewAST(AstVariable, "list"))

	closingType := OpeningClosingPairs[start.Type]
loop:
	for {
		token, err := p.scanner.Advance()
		if err != nil {
			return nil, err
		}
		switch token.Type {
		case closingType:
			ret.AddToken(token)
			break loop
		case TokInteger:
			ret.AddChild(NewASTAtom(AstInteger, token))
		case TokRational:
			ret.AddChild(NewASTAtom(AstRational, token))
		case TokReal:
			ret.AddChild(NewASTAtom(AstReal, token))
		case TokString:
			ret.AddChild(NewASTAtom(AstString, token))
		case TokIdentifier:
			symbol := NewAST(AstSymbol, "'"+token.Lexeme)
			symbol.AddToken(token)
			ret.AddChild(symbol)
		case TokSymbol:
			ret.AddChild(NewASTAtom(AstSymbol, token))
		case TokLeftParen, TokLeftBrace, TokLeftBracket:
			var quote *AST
			quote, err = p.parseQuote(token)
			if err != nil {
				return nil, err
			}
			ret.AddChild(quote)
		case TokQuotParen, TokQuotBrace, TokQuotBracket:
			return nil, errors.New("nested quotes are not allowed")
		default:
			ret.AddChild(NewASTAtom(AstQuoted, token))
		}
	}

	return ret, nil
}

func (p *Parser) parseDefinition(start *Token) (*AST, error) {
	closingType := OpeningClosingPairs[start.Type]

	tokenDef, err := p.expect(TokDef) // scan def keyword
	if err != nil {
		return nil, err
	}

	var token, identifier *Token
	token, err = p.expect(TokIdentifier, TokLeftParen, TokLeftBracket, TokLeftBrace)
	if err != nil {
		return nil, err
	}
	switch token.Type {
	case TokIdentifier:
		identifier = token
	default:
		return p.parseFuncDef(start, tokenDef, token)
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

func (p *Parser) parseStructDefinition(start *Token) (*AST, error) {
	structEndType := OpeningClosingPairs[start.Type]

	tokenStructDef, err := p.expect(TokDefStruct) // scan def keyword
	if err != nil {
		return nil, err
	}

	var identifier *Token
	identifier, err = p.expect(TokIdentifier)
	if err != nil {
		return nil, err
	}

	ret := NewAST(AstStructDefinition, identifier.Lexeme)
	ret.AddToken(start)
	ret.AddToken(tokenStructDef)
	ret.AddToken(identifier)

	var fieldsStart *Token
	fieldsStart, err = p.expect(TokLeftParen, TokLeftBrace, TokLeftBracket)
	if err != nil {
		return nil, err
	}
	ret.AddToken(fieldsStart)

	fieldsEndType := OpeningClosingPairs[fieldsStart.Type]

	var token *Token
	for {
		token, err = p.expect(TokIdentifier, fieldsEndType)
		if err != nil {
			return nil, err
		}
		if token.Type == TokIdentifier {
			ret.AddChild(NewASTAtom(AstStructField, token))
		} else {
			ret.AddToken(token)
			break
		}
	}

	structEnd, err := p.expect(structEndType)
	if err != nil {
		return nil, err
	}
	ret.AddToken(structEnd)

	return ret, nil
}

func (p *Parser) parseVarChange(start *Token) (*AST, error) {
	var setBang, identifier *Token
	var err error

	closingType := OpeningClosingPairs[start.Type]

	setBang, err = p.expect(TokSetBang) // scan def keyword
	if err != nil {
		return nil, err
	}

	identifier, err = p.expect(TokIdentifier)
	if err != nil {
		return nil, err
	}

	ret := NewAST(AstVarChange, identifier.Lexeme)
	ret.AddToken(start)
	ret.AddToken(setBang)
	ret.AddToken(identifier)

	var value *AST
	value, err = p.parseExpr(nil)
	if err != nil {
		return nil, err
	}
	ret.AddChild(value)

	var end *Token
	end, err = p.expect(closingType)
	if err != nil {
		return nil, err
	}
	ret.AddToken(end)

	return ret, nil
}

func (p *Parser) parseFuncDef(start, tokenDef, startFunc *Token) (*AST, error) {
	endType := OpeningClosingPairs[start.Type]
	endFuncType := OpeningClosingPairs[startFunc.Type]

	var identifier *Token
	var err error

	identifier, err = p.expect(TokIdentifier)
	if err != nil {
		return nil, err
	}

	ret := NewAST(AstDefinition, identifier.Lexeme)
	ret.AddToken(start)
	ret.AddToken(tokenDef)

	lambda := NewAST(AstLambda, identifier.Lexeme)
	lambda.AddToken(startFunc)
	lambda.AddToken(identifier)

	params := NewAST(AstParameters, "")
	var token, tokenEndParams *Token
	expectedTypes := []TokenType{TokVarParam, TokIdentifier, endFuncType}
paramsLoop:
	for {
		token, err = p.expect(expectedTypes...)
		if err != nil {
			return nil, err
		}
		switch token.Type {
		case TokIdentifier:
			params.AddChild(NewASTAtom(AstVariable, token))
		case TokVarParam:
			params.AddChild(NewASTAtom(AstVarParam, token))
			expectedTypes = expectedTypes[1:] // var param must be the last param
		case endFuncType:
			tokenEndParams = token
			break paramsLoop
		default:
			return nil, errors.New("expected identifier as parameter")
		}
	}

	lambda.AddChild(params)
	lambda.AddToken(tokenEndParams)

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

	lambda.AddChild(block)
	ret.AddChild(lambda)
	ret.AddToken(tokenEnd)

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
	expectedTypes := []TokenType{TokVarParam, TokIdentifier, endParamsType}
paramsLoop:
	for {
		token, err = p.expect(expectedTypes...)
		if err != nil {
			return nil, err
		}
		switch token.Type {
		case TokIdentifier:
			params.AddChild(NewASTAtom(AstVariable, token))
		case TokVarParam:
			params.AddChild(NewASTAtom(AstVarParam, token))
			expectedTypes = expectedTypes[1:]
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
