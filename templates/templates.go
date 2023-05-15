package templates

import (
	"embed"
	"fmt"
	"github.com/ystv/computing_site/link"
	"github.com/ystv/computing_site/team"
	"html/template"
	"io"
)

//go:embed *.tmpl
var tmpls embed.FS

type (
	Templater struct {
		dashboard *template.Template
	}
	DashboardParams struct {
		Link *link.Link
		Team *[]team.Member
	}
)

func (t *Templater) RenderTemplate(w io.Writer, data *DashboardParams, mainTmpl string) error {
	t1, err := template.New("base.tmpl").ParseFS(tmpls, "base.tmpl", mainTmpl)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return t1.Execute(w, data)
}
