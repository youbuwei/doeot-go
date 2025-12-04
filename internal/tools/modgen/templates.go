package modgen

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

var (
	domainTmpl   *template.Template
	appTmpl      *template.Template
	repoTmpl     *template.Template
	endpointTmpl *template.Template
	moduleTmpl   *template.Template
	modulesTmpl  *template.Template
)

func init() {
	domainTmpl = mustParse("domain.tmpl")
	appTmpl = mustParse("app.tmpl")
	repoTmpl = mustParse("repo.tmpl")
	endpointTmpl = mustParse("endpoint.tmpl")
	moduleTmpl = mustParse("module.tmpl")
	modulesTmpl = mustParse("modules.tmpl")
}

func mustParse(name string) *template.Template {
	t, err := template.ParseFS(templatesFS, "templates/"+name)
	if err != nil {
		panic(err)
	}
	return t
}
