package main

import (
	"fmt"
	"github.com/hlfstr/aeridya"
	"github.com/hlfstr/aeridya/page"
	"github.com/hlfstr/aeridya/router"
	"github.com/hlfstr/aeridya/router/staticrouter"
	"net/http"
	"os"
)

func main() {
	a, _, e := aeridya.Create("./conf")
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	r := staticrouter.Route{}
	r.Defaults()
	a.Route = r

	index, e := page.New("/", "Home", "home", "./here", "nothing")
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	index.Defaults()
	f := func(w http.ResponseWriter, r *http.Request, resp *router.Response) {
		resp.Status = 200
		fmt.Fprintf(w, "Hello from /\n")
		return
	}
	index.OnGet(f)
	r.Pages[index.Route] = index

	e = a.Run()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
