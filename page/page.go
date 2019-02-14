package page

import (
	"errors"
	"fmt"
	"github.com/hlfstr/aeridya/router"
	"html/template"
	"net/http"
	"net/http/httputil"
)

type Handler interface {
	//	Register()
	Run(w http.ResponseWriter, r *http.Request) *router.Response
}

type Page struct {
	Route string
	Title string

	templates []string
	Template  *template.Template

	options []string

	onGet    func(w http.ResponseWriter, r *http.Request, resp *router.Response)
	onPut    func(w http.ResponseWriter, r *http.Request, resp *router.Response)
	onPost   func(w http.ResponseWriter, r *http.Request, resp *router.Response)
	onDelete func(w http.ResponseWriter, r *http.Request, resp *router.Response)
}

func New(route string, title string, tmpldir string, tmpls ...string) (*Page, error) {
	p := &Page{}
	p.Route = route
	p.Title = title
	p.templates = addPath(tmpldir, tmpls)
	p.options = make([]string, 0)
	err := p.LoadPage()
	return p, err
}

func (p *Page) Defaults() {
	p.onGet = p.undefined
	p.onPut = p.undefined
	p.onPost = p.undefined
	p.onDelete = p.undefined
}

func (p *Page) LoadPage() error {
	var err error
	p.Template, err = template.ParseFiles(p.templates...)
	return err
}

func addPath(dir string, tmps []string) []string {
	s := make([]string, len(tmps))
	for i := 0; i < len(tmps); i++ {
		s[i] = dir + "/" + tmps[i]
	}
	return s
}

func (p *Page) OnGet(f func(w http.ResponseWriter, r *http.Request, resp *router.Response)) {
	p.options = append(p.options, "GET")
	p.onGet = f
}

func (p *Page) OnPut(f func(w http.ResponseWriter, r *http.Request, resp *router.Response)) {
	p.options = append(p.options, "PUT")
	p.onPut = f
}

func (p *Page) OnPost(f func(w http.ResponseWriter, r *http.Request, resp *router.Response)) {
	p.options = append(p.options, "POST")
	p.onPost = f
}

func (p *Page) OnDelete(f func(w http.ResponseWriter, r *http.Request, resp *router.Response)) {
	p.options = append(p.options, "DELETE")
	p.onDelete = f
}

func (p *Page) onHead(w http.ResponseWriter, r *http.Request, resp *router.Response) {
	requestDump, err := httputil.DumpRequest(r, false)
	if err != nil {
		resp.Status = 500
		resp.Error = err
		resp.Data = resp
		return
	}
	resp.Status = 200
	fmt.Fprintf(w, "%s\n", string(requestDump))
}

func (p *Page) onOptions(w http.ResponseWriter, r *http.Request, resp *router.Response) {
	resp.Status = 200
	fmt.Fprintf(w, "%s\n", p.options)
}

func (p *Page) undefined(w http.ResponseWriter, r *http.Request, resp *router.Response) {
	resp.Status = 400
	resp.Error = errors.New("Undefined Request " + r.Method + " to " + r.URL.Path)
	resp.Data = resp
}

func (p *Page) teapot(w http.ResponseWriter, r *http.Request, resp *router.Response) {
	resp.Status = 418
	resp.Error = errors.New("Unsupported Request " + r.Method + " to " + r.URL.Path)
	resp.Data = resp
}

func (p Page) Run(w http.ResponseWriter, r *http.Request) *router.Response {
	resp := &router.Response{}
	switch r.Method {
	case "GET":
		p.onGet(w, r, resp)
	case "PUT":
		p.onPut(w, r, resp)
	case "POST":
		p.onPost(w, r, resp)
	case "DELETE":
		p.onDelete(w, r, resp)
	case "OPTIONS":
		p.onOptions(w, r, resp)
	case "HEAD":
		p.onHead(w, r, resp)
	default:
		p.teapot(w, r, resp)
	}
	return resp
}
