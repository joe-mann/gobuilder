package executor

import (
	"io"

	"github.com/golang/go/src/pkg/text/template"
)

type Executor struct {
	funcs funcs
}

func New(relpath string) *Executor {
	e := &Executor{}
	e.funcs.Imports = []string{getBImport(), relpath}
	return e
}

func (e *Executor) AddBuildFile(r io.Reader, path string) error {
	funcs, err := findBuildFuncs(r, path)
	if err != nil {
		return err
	}
	e.funcs.Builders = append(e.funcs.Builders, funcs...)
	return nil
}

func (e *Executor) WriteBuildMain(w io.Writer) error {
	if err := mainTmpl.Execute(w, e.funcs); err != nil {
		return err
	}

	return nil
}

type funcs struct {
	Imports  []string
	Builders []builderFunc
}

var mainTmpl = template.Must(template.New("main").Parse(`
package main

import (
{{range .Imports}}
	"{{.}}"
{{end}}
)

var builders = []building.Builder{
{{range .Builders}}
	{Name: "{{.Name}}", F: {{.Package}}.{{.Name}}},
{{end}}
}

func main() {
	building.Main(builders)
}
`))
