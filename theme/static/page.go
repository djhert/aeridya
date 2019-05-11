package static

import (
	"github.com/hlfstr/aeridya"
	"html/template"
)

type Page struct {
	aeridya.Page

	Templates []string
	PageTemp  *template.Template
}

func (p *Page) PageInit(title, route, tmpldir string, tmpls ...string) error {
	p.Route = route
	p.Title = title
	p.Templates = AddDir(tmpldir, tmpls)
	err := p.LoadPage()
	return err
}

func AddDir(tmpldir string, tmpls []string) []string {
	t := make([]string, len(tmpls))
	for i := range tmpls {
		t[i] = tmpldir + "/" + tmpls[i]
	}
	return t
}

func (p *Page) LoadPage() error {
	var err error
	p.PageTemp, err = template.ParseFiles(p.Templates...)
	return err
}

func (p *Page) Get(resp *aeridya.Response) {
	resp.Good(200)
	if aeridya.Development {
		p.LoadPage()
	}
	p.PageTemp.Execute(resp.W, p)
}
