package interpreter

import (
	"testing"
)

func TestSequenceFunctions(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "take from nil",
			args: args{
				code: `(take 1 nil)`,
			},
			want: "nil",
		},
		{
			name: "take from list",
			args: args{
				code: `(take 2 (list 2 3 4))`,
			},
			want: "(list 2 3)",
		},
		{
			name: "take from vector",
			args: args{
				code: `(take 5 (vector 2 3 4))`,
			},
			want: "(vector 2 3 4)",
		},
		{
			name: "take from stream",
			args: args{
				code: `(take 2 (list->stream '(2 3 4)))`,
			},
			want: "(vector 2 3)",
		},
		{
			name: "take-while from list",
			args: args{
				code: `(take-while (λ (x) (< x 4)) (list 2 3 4))`,
			},
			want: "(list 2 3)",
		},
		{
			name: "take-while from vector",
			args: args{
				code: `(take-while (λ (x) (< x 4)) (vector 2 3 4))`,
			},
			want: "(vector 2 3)",
		},
		{
			name: "take-while from stream",
			args: args{
				code: `(take-while (λ (x) (< x 4)) (list->stream '(2 3 4)))`,
			},
			want: "(vector 2 3)",
		},
		{
			name: "filter a list",
			args: args{
				code: `
					(def (even? x) (= 0 (% x 2)))
					(def lst '(1 2 3 4 5 6))
					(filter even? lst)`,
			},
			want: "(list 2 4 6)",
		},
		{
			name: "filter a vector",
			args: args{
				code: `
					(def (even? x) (= 0 (% x 2)))
					(def v (vector 1 2 3 4 5 6))
					(filter even? v)`,
			},
			want: "(vector 2 4 6)",
		},
		{
			name: "filter a stream",
			args: args{
				code: `
					(def (even? x) (= 0 (% x 2)))
					(def s (list->stream '(1 2 3 4 5 6)))
					(take 3 (filter even? s))`,
			},
			want: "(vector 2 4 6)",
		},
		{
			name: "filter an iterator stream",
			args: args{
				code: `
					(def (even? x) (= 0 (% x 2)))
					(def s (iterator 1 (lambda (x) (+ x 1))))
					(take 3 (filter even? s))`,
			},
			want: "(vector 2 4 6)",
		},
		{
			name: "map a list",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def lst '(1 2 3 4 5 6))
					(map sq lst)`,
			},
			want: "(list 1 4 9 16 25 36)",
		},
		{
			name: "map a vector",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def v (vector 1 2 3 4 5 6))
					(map sq v)`,
			},
			want: "(vector 1 4 9 16 25 36)",
		},
		{
			name: "map a stream",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def s (list->stream (list 1 2 3 4 5 6)))
					(take 6 (map sq s))`,
			},
			want: "(vector 1 4 9 16 25 36)",
		},
		{
			name: "drop from a list",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def (even? x) (= 0 (% x 2)))
					(def lst '(1 2 3 4 5 6))
					(drop 1 (filter even? lst))`,
			},
			want: "(list 4 6)",
		},
		{
			name: "drop from a vector",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def (even? x) (= 0 (% x 2)))
					(def v (vector 1 2 3 4 5 6))
					(drop 1 (filter even? v))`,
			},
			want: "(vector 4 6)",
		},
		{
			name: "drop from a stream",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def (even? x) (= 0 (% x 2)))
					(def s (list->stream (list 1 2 3 4 5 6)))
					(take 10 (drop 1 (filter even? s)))`,
			},
			want: "(vector 4 6)",
		},
		{
			name: "drop-while from a list",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def (even? x) (= 0 (% x 2)))
					(def lst '(1 2 3 4 5 6))
					(drop-while (λ (x) (< x 4)) (filter even? lst))`,
			},
			want: "(list 4 6)",
		},
		{
			name: "drop-while from a vector",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def (even? x) (= 0 (% x 2)))
					(def v (vector 1 2 3 4 5 6))
					(drop-while (λ (x) (< x 4)) (filter even? v))`,
			},
			want: "(vector 4 6)",
		},
		{
			name: "drop-while from a stream",
			args: args{
				code: `
					(def (sq x) (* x x))
					(def (even? x) (= 0 (% x 2)))
					(def s (list->stream (list 1 2 3 4 5 6)))
					(take 10 (drop-while (λ (x) (< x 4)) (filter even? s)))`,
			},
			want: "(vector 4 6)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("Run() got = %s, want %s", got, tt.want)
			}
		})
	}

}
