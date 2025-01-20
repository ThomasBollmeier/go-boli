package frontend

import "testing"

func TestFindTailCalls(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "find tail calls",
			code: `
				(def (count-down n)
					(if (= n 0)
						nil
						(block
							(displayln n)
							(count-down (- n 1)))))`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast := createAst(tt.code)
			findTailCalls(ast)
			writeAst(ast, 0)
		})
	}
}

func createAst(code string) *AST {
	parser := NewParser()
	ast, err := parser.Parse(NewCharStreamString(code))
	if err != nil {
		panic(err)
	}
	return ast
}
