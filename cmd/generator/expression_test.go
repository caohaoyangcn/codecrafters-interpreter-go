package generator

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
			name: "ok",
			args: args{
				outDir: "./out",
				base:   "Expr",
				types: []string{
					"Binary   : Expr left, token.Token operator, Expr right",
					"Grouping : Expr expression",
					"Literal  : any value",
					"Unary    : token.Token operator, Expr right",
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
