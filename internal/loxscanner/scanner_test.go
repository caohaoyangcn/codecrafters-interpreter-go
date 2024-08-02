package loxscanner

import (
	"io"
	"strings"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestNewScanner(t *testing.T) {
	type args struct {
		src io.Reader
	}
	tests := []struct {
		name string
		args args
		want []*token.Token
	}{
		{
			name: "two-chars-token",
			args: args{
				strings.NewReader("!="),
			},
			want: []*token.Token{
				{
					Type:   token.BANG_EQUAL,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
		{
			name: "slash",
			args: args{
				strings.NewReader("/"),
			},
			want: []*token.Token{
				{
					Type:   token.SLASH,
					Lexeme: "",
					Object: nil,
					Line:   1,
				},
			},
		},
		{
			name: "comment-slash",
			args: args{
				strings.NewReader("//"),
			},
			want: []*token.Token{},
		},
		{
			name: "space",
			args: args{
				strings.NewReader(" "),
			},
			want: []*token.Token{},
		},
		{
			name: "space-2",
			args: args{
				strings.NewReader("\r"),
			},
			want: []*token.Token{},
		},
		{
			name: "space-3",
			args: args{
				strings.NewReader("\n"),
			},
			want: []*token.Token{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewScanner(tt.args.src)
			got.Scan()
			assert.Equal(t, tt.want, got.tokens)
		})
	}
}
