package main

import (
	"fmt"
	"github.com/hlfstr/aeridya"
	"net/http"
	"os"
)

func main() {
	e := aeridya.Create("./conf")
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	//	aeridya.Theme = &theme{}
	aeridya.Theme = &aeridya.Theming{}
	aeridya.Theme.Init()
	e = aeridya.Run()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

type theme struct {
	aeridya.Theming
}

func (t *theme) Init() {
}

func (t *theme) Serve(w http.ResponseWriter, r *http.Request, resp *aeridya.Response) {
	if r.URL.Path == "/" {
		resp.Status = 200
		fmt.Fprintf(w, "Hello Example!\n")
		return
	}
	t.Error(w, r, resp)
}

func (t *theme) Error(w http.ResponseWriter, r *http.Request, resp *aeridya.Response) {
	resp.Bad(404, "Example theme only supports /", w)
	fmt.Fprintf(w, "Error: %d\n%s\n", resp.Status, resp.Err)
	return
}
