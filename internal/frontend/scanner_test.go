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
		want    Token
		wantErr bool
	}{
		{
			name: "Get first token",
			fields: fields{
				stream: *NewBufferedStream(NewCharStreamString("   (+ 41 1)")),
			},
			want: *NewToken(TokLeftParen, "(", 1, 4),
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
