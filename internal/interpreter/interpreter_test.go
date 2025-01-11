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
			name: "evaluate power",
			args: args{
				code: "(^ 2,0 3,0 2)",
			},
			want: &Real{512.0},
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
