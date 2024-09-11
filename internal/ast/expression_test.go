package ast

import (
	"testing"
)

func TestDefineAst(t *testing.T) {
	type args struct {
		outDir string
		base   string
		types  []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "statement",
			args: args{
				outDir: ".",
				base:   "Stmt",
				types: []string{
					"Expression : Expr expression_",
					"Print      : Expr expression_",
				},
			},
		},
		{
			name: "ok",
			args: args{
				outDir: ".",
				base:   "Expr",
				types: []string{
					"Binary   : Expr left, token.Token operator, Expr right",
					"Grouping : Expr expression",
					"Literal  : any value",
					"Unary    : token.Token operator, Expr right",
					"Ternary  : Expr test, token.Token question, Expr left, token.Token colon, Expr right",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DefineAst(tt.args.outDir, tt.args.base, tt.args.types)
		})
	}
}
