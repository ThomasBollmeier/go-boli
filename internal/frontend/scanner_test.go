package frontend

import (
	"reflect"
	"testing"
)

func TestNewScanner(t *testing.T) {
	t.Run("Scanner created successfully", func(t *testing.T) {
		scanner := NewScanner(NewCharStreamString("(+ 41 1)"))
		if scanner == nil {
			t.Errorf("expected a Scanner, got nil")
		}
		if scanner.row != 1 {
			t.Errorf("expected row to be 1, got %d", scanner.row)
		}
		if scanner.col != 1 {
			t.Errorf("expected col to be 1, got %d", scanner.col)
		}
	})
}

func TestScanner_Advance(t *testing.T) {
	type fields struct {
		stream BufferedStream[rune]
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		{
			name: "Get first token",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("   (+ 41 1)")),
			},
			want: NewToken(TokLeftParen, "(", 1, 4),
		},
		{
			name: "Get nil",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("nil")),
			},
			want: NewToken(TokNil, "nil", 1, 1),
		},
		{
			name: "Get bool",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("#true")),
			},
			want: NewToken(TokBoolean, "#true", 1, 1),
		},
		{
			name: "Get bool",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("#t")),
			},
			want: NewToken(TokBoolean, "#t", 1, 1),
		},
		{
			name: "Get bool",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("#false")),
			},
			want: NewToken(TokBoolean, "#false", 1, 1),
		},
		{
			name: "Get bool",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("#f")),
			},
			want: NewToken(TokBoolean, "#f", 1, 1),
		},
		{
			name: "Get integer",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("42")),
			},
			want: NewToken(TokInteger, "42", 1, 1),
		},
		{
			name: "Get negative integer",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("-42")),
			},
			want: NewToken(TokInteger, "-42", 1, 1),
		},
		{
			name: "Get rational number",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("126/3")),
			},
			want: NewToken(TokRational, "126/3", 1, 1),
		},
		{
			name: "Get real number",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("+3,1415")),
			},
			want: NewToken(TokReal, "+3,1415", 1, 1),
		},
		{
			name: "Get a string",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("\"Thomas sagt: \\\"Hallo!\\\"\"")),
			},
			want: NewToken(TokString, "\"Thomas sagt: \\\"Hallo!\\\"\"", 1, 1),
		},
		{
			name: "Get a symbol",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("'i-am-a-symbol")),
			},
			want: NewToken(TokSymbol, "'i-am-a-symbol", 1, 1),
		},
		{
			name: "Get an identifier",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("#number-of-lines")),
			},
			want: NewToken(TokIdentifier, "#number-of-lines", 1, 1),
		},
		{
			name: "Get a keyword",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("\ndef")),
			},
			want: NewToken(TokDef, "def", 2, 1),
		},
		{
			name: "Get keyword and",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("\nand")),
			},
			want: NewToken(TokAnd, "and", 2, 1),
		},
		{
			name: "Get keyword or",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("or")),
			},
			want: NewToken(TokOr, "or", 1, 1),
		},
		{
			name: "Get spread",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("...xs")),
			},
			want: NewToken(TokSpread, "...xs", 1, 1),
		},
		{
			name: "Get var parameter",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("xs...")),
			},
			want: NewToken(TokVarParam, "xs...", 1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(&tt.fields.stream)
			got, err := s.Advance()
			if (err != nil) != tt.wantErr {
				t.Errorf("Advance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Advance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_AdvanceMany(t *testing.T) {
	assertCode(
		"(+ 41 1)",
		[]*Token{
			NewToken(TokLeftParen, "(", 1, 1),
			NewToken(TokPlus, "+", 1, 2),
			NewToken(TokInteger, "41", 1, 4),
			NewToken(TokInteger, "1", 1, 7),
			NewToken(TokRightParen, ")", 1, 8),
		},
		t,
	)
}

func TestScanner_AdvanceComparisonOperators(t *testing.T) {
	assertCode(
		"= > >= < <=",
		[]*Token{
			NewToken(TokEqual, "=", 1, 1),
			NewToken(TokGreater, ">", 1, 3),
			NewToken(TokGreaterEq, ">=", 1, 5),
			NewToken(TokLess, "<", 1, 8),
			NewToken(TokLessEq, "<=", 1, 10),
		},
		t,
	)
}

func TestScanner_LineComment(t *testing.T) {
	assertCode(
		`; a line comment
(def answer 42); <-- the answer to everything`,
		[]*Token{
			NewToken(TokLeftParen, "(", 2, 1),
			NewToken(TokDef, "def", 2, 2),
			NewToken(TokIdentifier, "answer", 2, 6),
			NewToken(TokInteger, "42", 2, 13),
			NewToken(TokRightParen, ")", 2, 15),
		},
		t,
	)
}

func TestScanner_BlockComment(t *testing.T) {
	assertCode(
		`#| 
This is a block comment 
|#
(def answer 42)`,
		[]*Token{
			NewToken(TokLeftParen, "(", 4, 1),
			NewToken(TokDef, "def", 4, 2),
			NewToken(TokIdentifier, "answer", 4, 6),
			NewToken(TokInteger, "42", 4, 13),
			NewToken(TokRightParen, ")", 4, 15),
		},
		t,
	)
}

func assertCode(code string, expectedTokens []*Token, t *testing.T) {
	stream := NewBufferedStream(NewCharStreamString(code))
	scanner := NewScanner(stream)

	actualTokens := getAllTokens(scanner)

	if len(actualTokens) != len(expectedTokens) {
		t.Errorf("expected %d tokens, got %d", len(expectedTokens), len(actualTokens))
		return
	}

	for i, actualToken := range actualTokens {
		expectedToken := expectedTokens[i]
		if !reflect.DeepEqual(*actualToken, *expectedToken) {
			t.Errorf("token %d: expected %v, got %v", i, *expectedToken, *actualToken)
		}
	}
}

func getAllTokens(scanner *Scanner) []*Token {
	var ret []*Token
	for {
		tok, err := scanner.Advance()
		if err != nil {
			break
		}
		ret = append(ret, tok)
	}

	return ret
}
