package bizgen

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

var (
	httpTmpl *template.Template
	rpcTmpl  *template.Template
)

func init() {
	var err error
	httpTmpl, err = template.ParseFS(templatesFS, "templates/http.tmpl")
	if err != nil {
		panic(err)
	}
	rpcTmpl, err = template.ParseFS(templatesFS, "templates/rpc.tmpl")
	if err != nil {
		panic(err)
	}
}
