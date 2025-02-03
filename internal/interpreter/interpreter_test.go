package interpreter

import (
	"go-boli/internal/frontend"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    ValueObject
		wantErr bool
	}{
		{
			name: "evaluate #true",
			args: args{
				code: "#true",
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate #t",
			args: args{
				code: "#t",
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate #false",
			args: args{
				code: "#false",
			},
			want: &Boolean{false},
		},
		{
			name: "evaluate #f",
			args: args{
				code: "#f",
			},
			want: &Boolean{false},
		},
		{
			name: "evaluate integer",
			args: args{
				code: "42",
			},
			want: &Integer{42},
		},
		{
			name: "evaluate addition",
			args: args{
				code: "(+ 41 1)",
			},
			want: &Integer{42},
		},
		{
			name: "evaluate subtraction",
			args: args{
				code: "(- 43 1)",
			},
			want: &Integer{42},
		},
		{
			name: "evaluate multiplication",
			args: args{
				code: "[* 6 7]",
			},
			want: &Integer{42},
		},
		{
			name: "evaluate division",
			args: args{
				code: "(/ 21 3)",
			},
			want: &Integer{7},
		},
		{
			name: "evaluate division",
			args: args{
				code: "(/ 1 3)",
			},
			want: &Rational{1, 3},
		},
		{
			name: "evaluate modulo",
			args: args{
				code: "(% 30 7)",
			},
			want: &Integer{2},
		},
		{
			name: "evaluate power",
			args: args{
				code: "(^ 2 3 2)",
			},
			want: &Integer{512},
		},
		{
			name: "evaluate rational",
			args: args{
				code: "2/3",
			},
			want: &Rational{2, 3},
		},
		{
			name: "evaluate rational as integer",
			args: args{
				code: "21/3",
			},
			want: &Integer{7},
		},
		{
			name: "evaluate rational addition",
			args: args{
				code: "(+ 1/3 2/3)",
			},
			want: &Integer{1},
		},
		{
			name: "evaluate rational subtraction",
			args: args{
				code: "(- 1/2 1/3)",
			},
			want: &Rational{1, 6},
		},
		{
			name: "evaluate rational multiplication",
			args: args{
				code: "(* 1/4 2/3)",
			},
			want: &Rational{1, 6},
		},
		{
			name: "evaluate rational division",
			args: args{
				code: "(/ 1/4 2/3)",
			},
			want: &Rational{3, 8},
		},
		{
			name: "evaluate real number",
			args: args{
				code: "42,0",
			},
			want: &Real{42},
		},
		{
			name: "evaluate real addition",
			args: args{
				code: "(+ 41,0 1)",
			},
			want: &Real{42.0},
		},
		{
			name: "evaluate real subtraction",
			args: args{
				code: "(- 44,0 4/2)",
			},
			want: &Real{42.0},
		},
		{
			name: "evaluate real multiplication",
			args: args{
				code: "[* 6,0 7,0]",
			},
			want: &Real{42.0},
		},
		{
			name: "evaluate real division",
			args: args{
				code: "(/ 21 3,0)",
			},
			want: &Real{7.0},
		},
		{
			name: "evaluate real power",
			args: args{
				code: "(^ 2,0 3,0 2)",
			},
			want: &Real{512.0},
		},
		{
			name: "evaluate string",
			args: args{
				code: "\"Thomas sagt: \\\"Hallo!\\\"\"",
			},
			want: &Str{"Thomas sagt: \"Hallo!\""},
		},
		{
			name: "evaluate conjunction",
			args: args{
				code: `
					(def num 5)
					(and (< 0 num) "Thomas")`,
			},
			want: &Str{"Thomas"},
		},
		{
			name: "evaluate disjunction",
			args: args{
				code: `
					(or nil "Thomas")`,
			},
			want: &Str{"Thomas"},
		},
		{
			name: "evaluate not",
			args: args{
				code: `(not nil)`,
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate definition",
			args: args{
				code: `(def answer 42)
						answer`,
			},
			want: &Integer{42},
		},
		{
			name: "value can be changed",
			args: args{
				code: `
					(def answer 41)
  					(set! answer 42)
					answer`,
			},
			want: &Integer{42},
		},
		{
			name: "evaluate equal",
			args: args{
				code: `
					(def answer (* 6 7))
					(= answer 42,0)`,
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate greater",
			args: args{
				code: `(> 2,0 3/2 1)`,
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate greater equal",
			args: args{
				code: `(>= 2,0 3/2 6/4 1)`,
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate less",
			args: args{
				code: `(< 1 3/2 2,0)`,
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate less equal",
			args: args{
				code: `(<= 1 1 3/2 2,0)`,
			},
			want: &Boolean{true},
		},
		{
			name: "evaluate if expression",
			args: args{
				code: `
					(def answer 126/3)
					(if (= answer 42)
    					"What is everything?"
						"Come ti chiami?")`,
			},
			want: &Str{"What is everything?"},
		},
		{
			name: "evaluate if expression 2",
			args: args{
				code: `
					(def answer 126/4)
					(if (= answer 42)
    					"What is everything?"
						"Wie geht's dir?")`,
			},
			want: &Str{"Wie geht's dir?"},
		},
		{
			name: "evaluate cond expression",
			args: args{
				code: `
					(def answer 126/4)
					(cond
						[(= answer 42) "What is everything?"]
						[#t "Wie geht's dir?"])`,
			},
			want: &Str{"Wie geht's dir?"},
		},
		{
			name: "evaluate cond expression without match",
			args: args{
				code: `
					(def answer 126/4)
					(cond
						[(= answer 42) "What is everything?"]
						[#f "Wie geht's dir?"])`,
			},
			want: &NilObject{},
		},
		{
			name: "evaluate block expression",
			args: args{
				code: `
					(def answer 43)
					(block
						(def answer 42)
						answer)`,
			},
			want: &Integer{42},
		},
		{
			name: "evaluate let expression",
			args: args{
				code: `
					(let ([forty 40]
                          [two 2]
                          [answer (+ forty two)])
						answer)`,
			},
			want: &Integer{42},
		},
		{
			name: "evaluate lambda call",
			args: args{
				code: `
					((λ (a b)
						(+ a b))
                      40 2)`,
			},
			want: &Integer{42},
		},
		{
			name: "evaluate function call",
			args: args{
				code: `
					(def (add a b)
            			(+ a b))
					(add 40 2)`,
			},
			want: &Integer{42},
		},
		{
			name: "evaluate closure call",
			args: args{
				code: `
					(def (make-adder n)
						(lambda (m) (+ m n)))
					((make-adder 2) 40)`,
			},
			want: &Integer{42},
		},
		{
			name: "evaluate with tail-call-optimization",
			args: args{
				code: `
					(def (sum n acc)
						(if (= n 0)
							acc
							(sum (- n 1) (+ acc n))))
					(sum 9 0)`,
			},
			want: &Integer{45},
		},
		{
			name: "evaluate varargs",
			args: args{
				code: `
					(def (my-add one two many...)
						(+ one two ...many))
					; (my-add 1 2 3 4) <-- ignored
					(my-add 1 2 3 4 5)`,
			},
			want: &Integer{15},
		},
		{
			name: "evaluate lambda with multiple arities",
			args: args{
				code: `
					(def (add a)
						(add a 42))
					(def (add a b)
						(+ a b))
				    (add 0)`,
			},
			want: &Integer{42},
		},
		{
			name: "evaluate pair",
			args: args{
				code: `
					(1 . (2 . nil))`,
			},
			want: &Pair{
				first: &Integer{1},
				second: &Pair{
					first:  &Integer{2},
					second: &NilObject{},
				},
			},
		},
		{
			name: "evaluate list function",
			args: args{
				code: `(list 1 2)`,
			},
			want: &Pair{
				first: &Integer{1},
				second: &Pair{
					first:  &Integer{2},
					second: &NilObject{},
				},
			},
		},
		{
			name: "evaluate list quote",
			args: args{
				code: `'(+ a b)`,
			},
			want: &Pair{
				first: &QuotedValue{
					token: &frontend.Token{
						Type:   frontend.TokPlus,
						Lexeme: "+",
						Row:    1,
						Col:    3,
					},
				},
				second: &Pair{
					first: &Symbol{"'a"},
					second: &Pair{
						first:  &Symbol{"'b"},
						second: &NilObject{},
					},
				},
			},
		},
		{
			name: "evaluate car",
			args: args{
				code: `
					(def pair (1 . (2 . nil)))
					(car pair)`,
			},
			want: &Integer{1},
		},
		{
			name: "evaluate cdr",
			args: args{
				code: `
					(def pair (1 . (2 . nil)))
					(cdr pair)`,
			},
			want: &Pair{
				first:  &Integer{2},
				second: &NilObject{},
			},
		},
		{
			name: "this pair is a list",
			args: args{
				code: `
					(def pair (1 . (2 . nil)))
					(list? pair)`,
			},
			want: &Boolean{true},
		},
		{
			name: "this pair is not a list",
			args: args{
				code: `
					(def pair (1 . 2))
					(list? pair)`,
			},
			want: &Boolean{false},
		},
		{
			name: "cons works",
			args: args{
				code: `
					(def pair (cons 1 2))
					(pair? pair)`,
			},
			want: &Boolean{true},
		},
		{
			name: "list-ref works",
			args: args{
				code: `
					(def lst (list 1 2 3))
					(list-ref lst 1)`,
			},
			want: &Integer{2},
		},
		{
			name: "vector-ref works",
			args: args{
				code: `
					(def v (vector 1 2 3))
					(vector-ref v 1)`,
			},
			want: &Integer{2},
		},
		{
			name: "vector? works",
			args: args{
				code: `
					(def v (vector 1 2 3))
					(vector? v)`,
			},
			want: &Boolean{true},
		},
		{
			name: "structure definition works",
			args: args{
				code: `
					(def-struct person (name first-name))
					person`,
			},
			want: &StructureType{
				name:   "person",
				fields: []string{"name", "first-name"},
			},
		},
		{
			name: "structure works",
			args: args{
				code: `
					(def-struct person (name first-name))
					(def ego (create-person "Bollmeier" "Thomas"))
					ego`,
			},
			want: &Structure{
				structType: &StructureType{
					name:   "person",
					fields: []string{"name", "first-name"},
				},
				values: map[string]ValueObject{
					"name":       &Str{"Bollmeier"},
					"first-name": &Str{"Thomas"},
				},
			},
		},
		{
			name: "structure getter works",
			args: args{
				code: `
				(def-struct person (name first-name))
				(def ego (create-person "Bollmeier" "Thomas"))
				(person-first-name ego)`,
			},
			want: &Str{"Thomas"},
		},
		{
			name: "structure setter works",
			args: args{
				code: `
				(def-struct person (name first-name))
				(def ego (create-person "Bollmeier" "Thomas"))
				(person-set-name! ego "Ballermeier")
				(person-name ego)`,
			},
			want: &Str{"Ballermeier"},
		},
		{
			name: "evaluate with shebang",
			args: args{
				code: `
				#! env boli

				(def (hello)
					"Hallo Welt!")
				(hello)`,
			},
			want: &Str{"Hallo Welt!"},
		},
		{
			name: "take works",
			args: args{
				code: `
				(def lst '(1 2 3 4))
				(take 1 (list->stream lst))`,
			},
			want: &Vector{elements: []ValueObject{&Integer{1}}},
		},
		{
			name: "filter works for stream",
			args: args{
				code: `
				(def (even? x) (= (% x 2) 0))
				(def lst '(1 2 3 4 5))
				(take 3 (filter even? (list->stream lst)))`,
			},
			want: &Vector{elements: []ValueObject{&Integer{2}, &Integer{4}}},
		},
		{
			name: "map works for stream",
			args: args{
				code: `
				(def (even? x) (= (% x 2) 0))
				(def (cube x) (* x x x))
				(def lst '(1 2 3 4 5))
				(take 3 (map cube (filter even? (list->stream lst))))`,
			},
			want: &Vector{elements: []ValueObject{&Integer{8}, &Integer{64}}},
		},
}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}
