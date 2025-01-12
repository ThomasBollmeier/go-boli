package frontend

import "testing"

func TestParser_ParseExpression(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name      string
		code      string
		wantType  AstType
		wantValue string
	}{
		{
			name:      "Parse integer",
			code:      "42",
			wantType:  AstInteger,
			wantValue: "42",
		},
		{
			name:      "Parse rational",
			code:      "  1/2  ",
			wantType:  AstRational,
			wantValue: "1/2",
		},
		{
			name:      "Parse real",
			code:      "\t3,1415\n",
			wantType:  AstReal,
			wantValue: "3,1415",
		},
		{
			name:      "Parse call",
			code:      "\t(+ 41 1)\n",
			wantType:  AstCall,
			wantValue: "",
		},
		{
			name:      "Parse call with braces",
			code:      "\t{+ 41 1}\n",
			wantType:  AstCall,
			wantValue: "",
		},
		{
			name:      "Parse definition",
			code:      "  (def answer 42)\n",
			wantType:  AstDefinition,
			wantValue: "answer",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program, err := parser.Parse(NewCharStreamString(test.code))
			if err != nil {
				t.Fatal(err)
			}
			ast := program.GetChildren()[0]
			if ast.astType != test.wantType {
				t.Errorf("got ast.astType %v, want %v", ast.astType, test.wantType)
			}
			if ast.value != test.wantValue {
				t.Errorf("got ast.value %v, want %v", ast.value, test.wantValue)
			}
		})
	}
}
