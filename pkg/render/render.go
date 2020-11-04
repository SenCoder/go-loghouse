package render

import (
	"github.com/Masterminds/sprig"
	"github.com/unrolled/render"
	"html/template"
)

var (
	App = render.New(render.Options{
		Directory:     "views",
		Extensions:    []string{".tmpl", ".html", "tpl"},
		IndentJSON:    true,
		Funcs:         []template.FuncMap{sprig.FuncMap()},
		IsDevelopment: true,
	})
)
