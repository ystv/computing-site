package templates

import (
	"embed"
	"html/template"
	"io"
	"log"
	"time"

	"github.com/ystv/computing_site/link"
	"github.com/ystv/computing_site/team"
)

//go:embed *.tmpl
var tmpls embed.FS

type (
	Templater struct{}

	DashboardParams struct {
		Link    *link.Link
		Team    *[]team.Member
		Commit  string
		Version string
	}
)

func (t *Templater) RenderTemplate(w io.Writer, data *DashboardParams, mainTmpl string) error {
	t1 := template.New("_base.tmpl")

	t1.Funcs(template.FuncMap{
		"thisYear": func() int {
			return time.Now().Year()
		},
	})

	t1, err := t1.ParseFS(tmpls, "_base.tmpl", "_top.tmpl", "_footer.tmpl", mainTmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	return t1.Execute(w, data)
}

// This section is for go template linter
var (
	AllTemplates = [][]string{
		{"dashboard.tmpl", "_base.tmpl", "_top.tmpl", "_footer.tmpl"},
	}

	_ = AllTemplates
)
