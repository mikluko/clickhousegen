package internal

import (
	_ "embed"
	"io"
	"sort"
	"text/template"
)

//go:embed render.go.tpl
var tpl string

type RenderData struct {
	PackageName string
	Struct      Struct
	Imports     []string
}

func Render(pkgName, structName string, obj Struct, w io.Writer) error {
	t := template.Must(template.New("").Parse(tpl))
	if structName != "" {
		obj.Name = structName
	}
	data := RenderData{
		PackageName: pkgName,
		Struct:      obj,
	}
	var imports []string
	for _, f := range obj.Fields {
		imports = append(imports, f.Imports...)
	}
	sort.Strings(imports)
	for i := range imports {
		if i > 0 && imports[i] == imports[i-1] {
			continue
		}
		data.Imports = append(data.Imports, imports[i])
	}
	return t.Execute(w, data)
}
