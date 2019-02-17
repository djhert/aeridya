package aeridya

import (
	"fmt"
	"net/http"
)

type theming interface {
	Serve(w http.ResponseWriter, r *http.Request, resp *Response)
	Error(w http.ResponseWriter, r *http.Request, resp *Response)
}

type ATheme struct {
	Page
}

func (t *ATheme) Init() {
	t.OnOptions("GET")
	t.Route = "/"
}

func (t *ATheme) Get(w http.ResponseWriter, r *http.Request, resp *Response) {
	resp.Good(200, w)
	fmt.Fprintf(w, "Hello Aeridya!\n")
	return
}

func (t *ATheme) Serve(w http.ResponseWriter, r *http.Request, resp *Response) {
	if r.URL.Path == "/" {
		ServePage(w, r, resp, t)
		return
	}
	t.Error(w, r, resp)
}

func (t *ATheme) Error(w http.ResponseWriter, r *http.Request, resp *Response) {
	resp.Bad(404, "Built-in theme only supports /", w)
	fmt.Fprintf(w, "Error: %d\n%s\n", resp.Status, resp.Err)
	return
}

func ThemeError(w http.ResponseWriter, r *http.Request, resp *Response, t theming) {
	t.Error(w, r, resp)
}
