package staticrouter

import (
	"errors"
	"github.com/hlfstr/aeridya/page"
	"github.com/hlfstr/aeridya/router"
	"html/template"
	"net/http"
)

type Route struct {
	Pages map[string]page.Handler
	Errs  map[int]*template.Template
}

func (s *Route) Defaults() {
	s.Errs = make(map[int]*template.Template)
	s.Pages = make(map[string]page.Handler)

	t := template.New("default")
	s.Errs[0], _ = t.Parse("An error has occurred: {{.Status}}\n{{.Error}}\n")

	tea := template.New("teapot")
	s.Errs[418], _ = tea.Parse("Error: {{.Status}}\nI am a teapot!\n")
}

// Parse removes trailing / from string if applicable
func (s Route) Parse(input string) (*router.Response, page.Handler) {
	r := &router.Response{}
	if len(input) > 1 {
		if input[len(input)-1:] == "/" {
			input = input[:len(input)-1]
		}
	}
	if val, ok := s.Pages[input]; ok {
		r.Status = 200
		return r, val
	}
	r.Status = 404
	r.Error = errors.New("Page not found " + input)
	return r, nil
}

func (s Route) Serve(w http.ResponseWriter, r *http.Request) *router.Response {
	resp, temp := s.Parse(r.URL.Path)
	if resp.Error != nil {
		return s.Error(resp, w, r)
	}
	resp = temp.Run(w, r)
	if resp.Error != nil {
		return s.Error(resp, w, r)
	}
	return resp
}

func (s Route) Error(resp *router.Response, w http.ResponseWriter, r *http.Request) *router.Response {
	if t, ok := s.Errs[resp.Status]; ok {
		t.Execute(w, resp.Data)
	} else {
		s.Errs[0].Execute(w, resp)
	}
	return resp
}
