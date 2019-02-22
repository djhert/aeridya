package aeridya

import (
	"fmt"
	"net/http/httputil"
)

type Paging interface {
	LoadPage() error
	Get(resp *Response)
	Put(resp *Response)
	Post(resp *Response)
	Delete(resp *Response)
	Options(resp *Response)
	Head(resp *Response)
	Unsupported(resp *Response)
}

type Page struct {
	Route string
	Title string

	options []string
}

func (p *Page) LoadPage() error {
	return nil
}

func (p *Page) Get(resp *Response) {
	p.undefined(resp)
}

func (p *Page) Put(resp *Response) {
	p.undefined(resp)
}

func (p *Page) Post(resp *Response) {
	p.undefined(resp)
}

func (p *Page) Delete(resp *Response) {
	p.undefined(resp)
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

func (p *Page) Head(resp *Response) {
	requestDump, err := httputil.DumpRequest(resp.R, false)
	if err != nil {
		resp.Status = 500
		resp.Err = err
		resp.Data = resp
		return
	}
	resp.Status = 200
	fmt.Fprintf(resp.W, "%s\n", string(requestDump))
}

func (p *Page) Options(resp *Response) {
	if p.options == nil {
		p.undefined(resp)
		return
	}
	resp.Good(200)
	fmt.Fprintf(resp.W, "%s\n", p.options)
}

func (p *Page) undefined(resp *Response) {
	resp.Bad(400, "Bad Request "+resp.R.Method+" to "+resp.R.URL.Path)
	fmt.Fprintf(resp.W, "Error: %d\n%s\n", resp.Status, resp.Err)
	resp.Data = resp
}

func (p *Page) Unsupported(resp *Response) {
	resp.Bad(418, "Unsupported Request "+resp.R.Method+" to "+resp.R.URL.Path)
	fmt.Fprintf(resp.W, "Error: %d\n%s\nI'M A TEAPOT!\n", resp.Status, resp.Err)
	resp.Data = resp
}

func ServePage(resp *Response, p Paging) {
	switch resp.R.Method {
	case "GET":
		p.Get(resp)
	case "PUT":
		p.Put(resp)
	case "POST":
		p.Post(resp)
	case "DELETE":
		p.Delete(resp)
	case "OPTIONS":
		p.Options(resp)
	case "HEAD":
		p.Head(resp)
	default:
		p.Unsupported(resp)
	}
}
