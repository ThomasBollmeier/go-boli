package frontend

import (
	"fmt"
	"testing"
)

func TestParser_ParseExpression(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name      string
		code      string
		wantType  AstType
		wantValue string
	}{
		{
			name:      "Parse nil",
			code:      "nil",
			wantType:  AstNil,
			wantValue: "nil",
		},
		{
			name:      "Parse boolean",
			code:      "#t",
			wantType:  AstBoolean,
			wantValue: "#t",
		},
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
			name:      "Parse string",
			code:      "\t\"3,1415\"\n",
			wantType:  AstString,
			wantValue: "\"3,1415\"",
		},
		{
			name:      "Parse symbol",
			code:      "\t'a-symbol",
			wantType:  AstSymbol,
			wantValue: "'a-symbol",
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
			name:      "Parse conjunction",
			code:      "(and (< 0 num) (< num 10))",
			wantType:  AstConjunction,
			wantValue: "",
		},
		{
			name:      "Parse disjunction",
			code:      "(or (<= num 0) (>= num 10))",
			wantType:  AstDisjunction,
			wantValue: "",
		},
		{
			name:      "Parse definition",
			code:      "  (def answer 42)\n",
			wantType:  AstDefinition,
			wantValue: "answer",
		},
		{
			name: "Parse if expression",
			code: `
					(if (= answer 42)
    					"What is everything?"
						"Come ti chiami?")`,
			wantType:  AstIfExpression,
			wantValue: "",
		},
		{
			name: "Parse cond expression",
			code: `
					(cond
						[(= answer 42) "What is everything?"]
						[#t	"Come ti chiami?"])`,
			wantType:  AstIfExpression,
			wantValue: "",
		},
		{
			name: "Parse block expression",
			code: `
				(block
					(def answer 42)
					answer)`,
			wantType:  AstBlock,
			wantValue: "",
		},
		{
			name: "Parse let expression",
			code: `
				(let ([forty 40]
                	  [two 2]
                      [answer (+ forty two)])
					answer)`,
			wantType:  AstBlock,
			wantValue: "",
		},
		{
			name: "Parse lambda expression",
			code: `
				(def fib (λ (n)
					     	(cond 
								[(= n 0) 0]
								[(= n 1) 1]
								[#t (+ (fib (- n 2)) (fib (- n 1)))])))`,
			wantType:  AstDefinition,
			wantValue: "fib",
		},
		{
			name: "Parse function definition",
			code: `
				(def (fib n)
					(cond 
						[(= n 0) 0]
						[(= n 1) 1]
						[#t (+ (fib (- n 2)) (fib (- n 1)))]))`,
			wantType:  AstDefinition,
			wantValue: "fib",
		},
		{
			name: "Parse pair",
			code: `
				(a . (b . nil))`,
			wantType:  AstPair,
			wantValue: "",
		},
		{
			name: "Parse set!",
			code: `
				(set! answer 42)`,
			wantType:  AstVarChange,
			wantValue: "answer",
		},
		{
			name: "Parse set!",
			code: `
				(set! answer 42)`,
			wantType:  AstVarChange,
			wantValue: "answer",
		},
		{
			name: "Parse varargs",
			code: `
				(def (my-add xs...)
					(+ ...xs))
				(my-add 1 2 3 4)`,
			wantType:  AstDefinition,
			wantValue: "my-add",
		},
		{
			name: "Parse varargs 2",
			code: `
				(def my-add (λ (xs...)
					(+ ...xs)))
				(my-add 1 2 3 4)`,
			wantType:  AstDefinition,
			wantValue: "my-add",
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
			writeAst(ast, 0)
			fmt.Println(ast.GetLexemes())
		})
	}
}

func writeAst(ast *AST, indent int) {
	const delta int = 2
	writeln("{", indent)
	writeln(fmt.Sprintf("type: %s", ast.astType), indent+delta)
	writeln(fmt.Sprintf("value: %s", ast.value), indent+delta)
	children := ast.GetChildren()
	if len(children) > 0 {
		writeln("children: [", indent+delta)
		for _, child := range children {
			writeAst(child, indent+2*delta)
		}
		writeln("]", indent+delta)
	}
	attrs := ast.GetAttributes()
	if len(attrs) > 0 {
		writeln("attributes: [", indent+delta)
		offset := indent + 2*delta
		for key, value := range attrs {
			writeln(fmt.Sprintf("%s: %v", key, value), offset)
		}
		writeln("]", indent+delta)
	}
	writeln("}", indent)
}

func write(text string, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
	fmt.Print(text)
}

func writeln(text string, indent int) {
	write(text+"\n", indent)
}
