package aeridya

import (
	"fmt"
)

type theming interface {
	Serve(resp *Response)
	Error(resp *Response)
}

type ATheme struct {
	Page
}

func (t *ATheme) Init() {
	t.OnOptions("GET")
	t.Route = "/"
}

func (t *ATheme) Get(resp *Response) {
	resp.Good(200)
	fmt.Fprintf(resp.W, "Hello Aeridya!\n")
	return
}

func (t *ATheme) Serve(resp *Response) {
	if resp.R.URL.Path == "/" {
		ServePage(resp, t)
		return
	}
	t.Error(resp)
}

func (t *ATheme) Error(resp *Response) {
	resp.Bad(404, "Built-in theme only supports /")
	fmt.Fprintf(resp.W, "Error: %d\n%s\n", resp.Status, resp.Err)
	return
}

func ThemeError(resp *Response, t theming) {
	t.Error(resp)
}
