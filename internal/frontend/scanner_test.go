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
			name: "Get integer",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("42")),
			},
			want: NewToken(TokInteger, "42", 1, 1),
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
				stream: *NewBufferedStream(NewCharStreamString("3,1415")),
			},
			want: NewToken(TokReal, "3,1415", 1, 1),
		},
		{
			name: "Get a string",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("\"Thomas sagt: \\\"Hallo!\\\"\"")),
			},
			want: NewToken(TokString, "\"Thomas sagt: \\\"Hallo!\\\"\"", 1, 1),
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
	code := "(+ 41 1)"
	stream := NewBufferedStream(NewCharStreamString(code))
	scanner := NewScanner(stream)

	expectedTokens := []*Token{
		NewToken(TokLeftParen, "(", 1, 1),
		NewToken(TokPlus, "+", 1, 2),
		NewToken(TokInteger, "41", 1, 4),
		NewToken(TokInteger, "1", 1, 7),
		NewToken(TokRightParen, ")", 1, 8),
	}

	actualTokens := getAllTokens(scanner)

	if len(actualTokens) != len(expectedTokens) {
		t.Errorf("expected %d tokens, got %d", len(expectedTokens), len(actualTokens))
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
