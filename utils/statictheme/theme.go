package statictheme

import (
	"github.com/hlfstr/aeridya"
	"html/template"
)

type Theme struct {
	Pages  map[string]aeridya.Paging
	Errors map[int]*template.Template

	TemplateDir string
}

func (t *Theme) StaticInit(base string) error {
	t.Pages = make(map[string]aeridya.Paging)
	t.Errors = make(map[int]*template.Template)

	a := template.New("default")
	t.Errors[0], _ = a.Parse("An error has occurred: {{.Status}}\n{{.Err}}\n")

	if s, err := aeridya.Config.GetString("", "Template"); err != nil {
		return err
	} else {
		t.TemplateDir = s
	}
	return nil
}

func (t Theme) Parse(input string, resp *aeridya.Response) aeridya.Paging {
	if p, ok := t.Pages[input]; !ok {
		resp.Error("Page " + input + " not found")
		return nil
	} else {
		return p
	}
}

func (t Theme) Serve(resp *aeridya.Response) {
	o := t.Parse(resp.R.URL.Path, resp)
	if o == nil {
		resp.Bad(404, resp.Err.Error())
		aeridya.ThemeError(resp, t)
		return
	}
	aeridya.ServePage(resp, o)
	return
}

func (t Theme) Error(resp *aeridya.Response) {
	if resp.Data == nil {
		resp.Data = resp
	}
	if s, ok := t.Errors[resp.Status]; ok {
		s.Execute(resp.W, resp.Data)
	} else {
		t.Errors[0].Execute(resp.W, resp)
	}
	return
}
