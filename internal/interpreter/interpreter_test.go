package interpreter

import (
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
			name: "evaluate definition",
			args: args{
				code: `(def answer 42)
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
