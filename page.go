package aeridya

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Paging interface {
	LoadPage() error
	Get(w http.ResponseWriter, r *http.Request, resp *Response)
	Put(w http.ResponseWriter, r *http.Request, resp *Response)
	Post(w http.ResponseWriter, r *http.Request, resp *Response)
	Delete(w http.ResponseWriter, r *http.Request, resp *Response)
	Options(w http.ResponseWriter, r *http.Request, resp *Response)
	Head(w http.ResponseWriter, r *http.Request, resp *Response)
	Unsupported(w http.ResponseWriter, r *http.Request, resp *Response)
}

type Page struct {
	Route string
	Title string

	options []string
}

func (p *Page) LoadPage() error {
	return nil
}

func (p *Page) Get(w http.ResponseWriter, r *http.Request, resp *Response) {
	p.undefined(w, r, resp)
}

func (p *Page) Put(w http.ResponseWriter, r *http.Request, resp *Response) {
	p.undefined(w, r, resp)
}

func (p *Page) Post(w http.ResponseWriter, r *http.Request, resp *Response) {
	p.undefined(w, r, resp)
}

func (p *Page) Delete(w http.ResponseWriter, r *http.Request, resp *Response) {
	p.undefined(w, r, resp)
}

func (p *Page) OnOptions(opts ...string) {
	p.options = make([]string, len(opts)+2)
	for i := range p.options {
		if i < len(opts) {
			p.options[i] = opts[i]
		} else {
			p.options[i] = "HEAD"
			p.options[i+1] = "OPTIONS"
			return
		}
	}
}

func (p *Page) Head(w http.ResponseWriter, r *http.Request, resp *Response) {
	requestDump, err := httputil.DumpRequest(r, false)
	if err != nil {
		resp.Status = 500
		resp.Err = err
		resp.Data = resp
		return
	}
	resp.Status = 200
	fmt.Fprintf(w, "%s\n", string(requestDump))
}

func (p *Page) Options(w http.ResponseWriter, r *http.Request, resp *Response) {
	if p.options == nil {
		p.undefined(w, r, resp)
		return
	}
	resp.Good(200, w)
	fmt.Fprintf(w, "%s\n", p.options)
}

func (p *Page) undefined(w http.ResponseWriter, r *http.Request, resp *Response) {
	resp.Bad(400, "Bad Request "+r.Method+" to "+r.URL.Path, w)
	fmt.Fprintf(w, "Error: %d\n%s\n", resp.Status, resp.Err)
	resp.Data = resp
}

func (p *Page) Unsupported(w http.ResponseWriter, r *http.Request, resp *Response) {
	resp.Bad(418, "Unsupported Request "+r.Method+" to "+r.URL.Path, w)
	fmt.Fprintf(w, "Error: %d\n%s\nI'M A TEAPOT!\n", resp.Status, resp.Err)
	resp.Data = resp
}

func ServePage(w http.ResponseWriter, r *http.Request, resp *Response, p Paging) {
	switch r.Method {
	case "GET":
		p.Get(w, r, resp)
	case "PUT":
		p.Put(w, r, resp)
	case "POST":
		p.Post(w, r, resp)
	case "DELETE":
		p.Delete(w, r, resp)
	case "OPTIONS":
		p.Options(w, r, resp)
	case "HEAD":
		p.Head(w, r, resp)
	default:
		p.Unsupported(w, r, resp)
	}
}
