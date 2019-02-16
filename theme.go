package aeridya

import (
	"fmt"
	"net/http"
)

type theming interface {
	Init()
	Serve(w http.ResponseWriter, r *http.Request, resp *Response)
	Error(w http.ResponseWriter, r *http.Request, resp *Response)
}

type Theming struct {
	Page
}

func (t *Theming) Init() {
	t.OnOptions("GET")
	t.Route = "/"
}

func (t *Theming) Get(w http.ResponseWriter, r *http.Request, resp *Response) {
	resp.Good(200, w)
	fmt.Fprintf(w, "Hello Aeridya!\n")
	return
}

func (t *Theming) Serve(w http.ResponseWriter, r *http.Request, resp *Response) {
	if r.URL.Path == "/" {
		ServePage(w, r, resp, t)
		return
	}
	t.Error(w, r, resp)
}

func (t *Theming) Error(w http.ResponseWriter, r *http.Request, resp *Response) {
	resp.Bad(404, "Built-in theme only supports /", w)
	fmt.Fprintf(w, "Error: %d\n%s\n", resp.Status, resp.Err)
	return
}
