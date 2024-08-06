package ast

import (
	"fmt"
	"go/format"
	"io"
	"os"
	"path"
	"strings"
)

func DefineAst(outDir string, base string, types []string) {
	// create dir if not exists
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.Mkdir(outDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	p := path.Join(outDir, base+".go")
	fd, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	sb := &strings.Builder{}
	sb.WriteString("package ast\n\n")
	sb.WriteString(`import "github.com/codecrafters-io/interpreter-starter-go/internal/token"

`)

	sb.WriteString(fmt.Sprintf(`type %s interface{
	Accept(visitor Visitor[any]) (any, error)
}

`, base))

	for _, t := range types {
		typeName := strings.Split(t, ":")[0]
		fields := strings.Split(t, ":")[1]

		constructorFunc := &strings.Builder{}
		constructorFuncReturn := &strings.Builder{}
		constructorFunc.WriteString(fmt.Sprintf(`func New%s%s(`, base, typeName))

		// struct definition
		sb.WriteString(fmt.Sprintf("type %s struct {\n", typeName))

		for _, s := range strings.Split(fields, ",") {
			s = strings.TrimSpace(s)
			fieldType := strings.Split(s, " ")[0]
			fieldName := strings.Title(strings.Split(s, " ")[1])
			sb.WriteString(fmt.Sprintf("\t%s\t%s\n", fieldName, fieldType))
			// constructor func params
			constructorFunc.WriteString(fmt.Sprintf("%s %s,", fieldName, fieldType))
			constructorFuncReturn.WriteString(fmt.Sprintf("%s: %s,", fieldName, fieldName))
		}
		sb.WriteString("}\n\n")
		constructorFunc.WriteString(fmt.Sprintf(") %s {", base))
		constructorFunc.WriteString(fmt.Sprintf("return &%s{", typeName))
		constructorFunc.WriteString(constructorFuncReturn.String())
		constructorFunc.WriteString("}}\n\n")

		sb.WriteString(constructorFunc.String())

		// visitor pattern
		sb.WriteString(fmt.Sprintf(`func (%s *%s) Accept(visitor Visitor[any]) (any, error) {
	return visitor.Visit%s%s(%s)
}

`,
			strings.ToLower(string(typeName[0])), typeName, base, typeName, strings.ToLower(string(typeName[0]))))
	}

	DefineVisitor(sb, base, types)

	formatted, err := format.Source([]byte(sb.String()))
	if err != nil {
		panic(err)
	}
	fd.Write(formatted)

}

func DefineVisitor(outDir io.StringWriter, base string, types []string) {
	outDir.WriteString("type Visitor[T any] interface { \n\n")
	for _, s := range types {
		typeName := strings.Split(s, ":")[0]
		outDir.WriteString(fmt.Sprintf("\tVisit%s%s(%s *%s) (T, error)\n",
			base, typeName, strings.ToLower(base), typeName))
	}
	outDir.WriteString("}\n\n")
}
