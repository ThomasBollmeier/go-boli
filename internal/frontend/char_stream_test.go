package frontend

import (
	"reflect"
	"testing"
)

func TestCharStreamString_Advance(t *testing.T) {
	type fields struct {
		chars []rune
		idx   int
	}
	tests := []struct {
		name    string
		fields  fields
		want    rune
		wantErr bool
	}{
		{
			name: "Advance to next char",
			fields: fields{
				chars: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
				idx:   1,
			},
			want:    'b',
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := &CharStreamString{
				chars: tt.fields.chars,
				idx:   tt.fields.idx,
			}
			got, err := stream.Advance()
			if (err != nil) != tt.wantErr {
				t.Errorf("Advance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Advance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCharStreamString(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want *CharStreamString
	}{
		{
			name: "NewCharStreamString works",
			args: args{
				text: "Über",
			},
			want: &CharStreamString{
				chars: []rune{'Ü', 'b', 'e', 'r'},
				idx:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCharStreamString(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCharStreamString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBufferedCharStream_Advance(t *testing.T) {
	type fields struct {
		stream CharStream
		buf    []rune
	}
	tests := []struct {
		name    string
		fields  fields
		want    rune
		wantErr bool
	}{
		{
			name: "Advance to next char without buffer",
			fields: fields{
				stream: NewCharStreamString("abcd"),
				buf:    []rune{},
			},
			want:    'a',
			wantErr: false,
		},
		{
			name: "Advance to next char with buffer",
			fields: fields{
				stream: NewCharStreamString("bcd"),
				buf:    []rune{'a'},
			},
			want:    'a',
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufCharStream := &BufferedCharStream{
				stream: tt.fields.stream,
				buf:    tt.fields.buf,
			}
			got, err := bufCharStream.Advance()
			if (err != nil) != tt.wantErr {
				t.Errorf("Advance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Advance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBufferedCharStream_Peek(t *testing.T) {
	type fields struct {
		stream CharStream
		buf    []rune
	}
	tests := []struct {
		name    string
		fields  fields
		want    rune
		wantErr bool
	}{
		{
			name: "Peek without buffer",
			fields: fields{
				stream: NewCharStreamString("abcd"),
				buf:    []rune{},
			},
			want:    'a',
			wantErr: false,
		},
		{
			name: "Peek with buffer",
			fields: fields{
				stream: NewCharStreamString("bcd"),
				buf:    []rune{'a'},
			},
			want:    'a',
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufCharStream := &BufferedCharStream{
				stream: tt.fields.stream,
				buf:    tt.fields.buf,
			}
			got, err := bufCharStream.Peek()
			if (err != nil) != tt.wantErr {
				t.Errorf("Peek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Peek() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBufferedCharStream_PeekMany(t *testing.T) {
	type fields struct {
		stream CharStream
		buf    []rune
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []rune
	}{
		{
			name: "PeekMany without buffer",
			fields: fields{
				stream: NewCharStreamString("abcd"),
				buf:    []rune{},
			},
			args: args{
				n: 10,
			},
			want: []rune{'a', 'b', 'c', 'd'},
		},
		{
			name: "PeekMany with buffer",
			fields: fields{
				stream: NewCharStreamString("cd"),
				buf:    []rune{'a', 'b'},
			},
			args: args{
				n: 10,
			},
			want: []rune{'a', 'b', 'c', 'd'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufCharStream := &BufferedCharStream{
				stream: tt.fields.stream,
				buf:    tt.fields.buf,
			}
			if got := bufCharStream.PeekMany(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeekMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBufferedCharStream(t *testing.T) {
	type args struct {
		stream CharStream
	}
	tests := []struct {
		name string
		args args
		want *BufferedCharStream
	}{
		{
			name: "Creation of buffered inStream works",
			args: args{
				stream: NewCharStreamString("abcd"),
			},
			want: &BufferedCharStream{
				stream: NewCharStreamString("abcd"),
				buf:    []rune{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBufferedStream(tt.args.stream); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBufferedCharStream() = %v, want %v", got, tt.want)
			}
		})
	}
}
