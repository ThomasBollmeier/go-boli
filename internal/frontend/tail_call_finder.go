package frontend

func findTailCalls(ast *AST) {
	children := ast.GetChildren()

	if ast.GetType() == AstLambda {
		body := children[1]
		setTailCall(body)
	}

	for _, child := range children {
		findTailCalls(child)
	}
}

func setTailCall(ast *AST) {
	astType := ast.GetType()
	switch astType {
	case AstBlock:
		children := ast.GetChildren()
		if len(children) > 0 {
			lastChild := children[len(children)-1]
			setTailCall(lastChild)
		}
	case AstIfExpression:
		children := ast.GetChildren()
		setTailCall(children[1]) // set tail call in consequent
		setTailCall(children[2]) // set tail call in alternative
	case AstCall:
		ast.SetAttribute(AstAttrKeyIsTailCall, true)
	}
}
