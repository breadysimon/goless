package annotations

import (
	"fmt"
	"go/ast"
	"go/token"

	"go/parser"
)

func applyAnnotations(path string) {
	for _, f := range listFiles(path) {
		if ans := findAnnotations(f); len(ans) > 0 {
			for _, an := range ans {
				g := CreateGenerator(&an)
				g.Run()
			}
		}
	}
}
func listFiles(path string) (fl []string) {
	return
}

type Field struct {
	Type string
	Tag  string
}

type Annotation struct {
	Command   string
	Arguments map[string]string
	Package   string
	Struct    map[string]Field
}

func findAnnotations(file string) (ans []Annotation) {
	return
}

type CodeGenerator interface {
	Run()
}
type AnnoApi struct {
}

func (o *AnnoApi) Run() {

}
func CreateGenerator(an *Annotation) (g CodeGenerator) {
	return &AnnoApi{}
}

func parse(file string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	// f, err := parser.ParseExpr(`AAa(xx 1, b "123")`)
	if err != nil {
		panic(err)
	}

	ast.Print(fset, f)
	ast.Inspect(f, func(n ast.Node) bool {
		if r, ok := n.(*ast.GenDecl); ok {
			c := r.Doc.List[0].Text
			if c[:len("// @")] == "// @" {
				x := r.Specs[0].(*ast.TypeSpec)
				t := x.Type.(*ast.StructType)
				ff := []map[string]string{}
				for _, v := range t.Fields.List {
					kv := make(map[string]string)
					kv["name"] = v.Names[0].Name
					kv["type"] = v.Type.(*ast.Ident).Name
					kv["note"] = v.Tag.Value
					ff = append(ff, kv)
				}
				fmt.Print(c, "\n", x.Name.Name, ff)

			}
		}
		return true
	})
}
