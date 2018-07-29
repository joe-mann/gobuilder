package executor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"reflect"
	"strings"

	"github.com/joe-mann/gobuilder/building"
)

type builderFunc struct {
	Name    string
	Package string
}

func findBuildFuncs(r io.Reader, path string) ([]builderFunc, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", r, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	pkg := f.Name.String()

	funcs := []builderFunc{}
	imports := make(map[string]string)

	var errOut error
	ast.Inspect(f, func(n ast.Node) bool {
		switch _n := n.(type) {
		case *ast.FuncDecl:
			name, ok, err := buildFunc(_n, path, fset.Position(n.Pos()), imports)
			if err != nil {
				errOut = err
			}
			if ok {
				funcs = append(funcs, builderFunc{Name: name, Package: pkg})
			}
		case *ast.ImportSpec:
			name, path := getImport(_n)
			imports[name] = path
		}
		return errOut == nil
	})
	return funcs, errOut
}

func getImport(i *ast.ImportSpec) (string, string) {
	if i.Path == nil {
		return "", ""
	}
	if i.Name != nil {
		return strings.Trim(i.Name.Name, "\""), strings.Trim(i.Path.Value, "\"")
	}
	pos := strings.LastIndex(i.Path.Value, "/")
	if pos == -1 {
		return strings.Trim(i.Path.Value, "\""), strings.Trim(i.Path.Value, "\"")
	}
	return strings.Trim(i.Path.Value[pos+1:len(i.Path.Value)], "\""), strings.Trim(i.Path.Value, "\"")
}

func buildFunc(f *ast.FuncDecl, path string, pos token.Position, imports map[string]string) (string, bool, error) {
	name := f.Name.Name
	if !strings.HasPrefix(name, "Build") {
		return "", false, nil
	}
	if f.Type.Params.NumFields() != 1 {
		return "", false, wrongSignature(name, path, pos)
	}
	if !checkParam(f.Type.Params.List[0], imports) {
		return "", false, wrongSignature(name, path, pos)
	}
	if f.Type.Results.NumFields() != 0 {
		return "", false, wrongSignature(name, path, pos)
	}
	return name, true, nil
}

func checkParam(param *ast.Field, imports map[string]string) bool {
	star, ok := param.Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	importPath, ok := imports[ident.Name]
	if !ok || importPath != getBImport() {
		return false
	}
	return true
}

func getBImport() string {
	b := building.B{}
	return reflect.TypeOf(b).PkgPath()
}

func wrongSignature(name string, path string, pos token.Position) error {
	return fmt.Errorf("%s:%s: wrong signature for %s, must be func %s(b *%s.B)", path, pos, name, name, getBImport())
}
