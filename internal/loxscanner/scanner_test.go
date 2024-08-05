package loxscanner

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestNewScanner(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want []*token.Token
	}{
		{
			name: "two-chars-token",
			args: args{
				"!=",
			},
			want: []*token.Token{
				{
					Type:   token.BANG_EQUAL,
					Lexeme: token.BANG_EQUAL.Repr(nil),
					Object: nil,
					Line:   1,
				},
				{
					Type:   token.EOF,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
		{
			name: "slash",
			args: args{
				"/",
			},
			want: []*token.Token{
				{
					Type:   token.SLASH,
					Lexeme: token.SLASH.Repr(nil),
					Object: nil,
					Line:   1,
				}, {
					Type:   token.EOF,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
		{
			name: "comment-slash",
			args: args{
				"//Comment",
			},
			want: []*token.Token{
				{
					Type:   token.EOF,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
		{
			name: "space",
			args: args{
				" ",
			},
			want: []*token.Token{
				{
					Type:   token.EOF,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
		{
			name: "space-2",
			args: args{
				"\r",
			},
			want: []*token.Token{
				{
					Type:   token.EOF,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
		{
			name: "space-3",
			args: args{
				"\n",
			},
			want: []*token.Token{
				{
					Type:   token.EOF,
					Lexeme: "",
					Object: nil,
					Line:   2,
				},
			},
		},
		{
			name: "string",
			args: args{
				"\"hello\"",
			},
			want: []*token.Token{
				{
					Type:   token.STRING,
					Lexeme: "\"hello\"",
					Object: "hello",
					Line:   1,
				},
				{
					Type:   token.EOF,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewScanner(tt.args.src)
			got.ScanAll()
			assert.Equal(t, tt.want, got.tokens)
		})
	}
}
