package main

import (
	"fmt"
	"github.com/hlfstr/aeridya"
	"github.com/hlfstr/aeridya/utils/statictheme"
	"html/template"
	"os"
)

func main() {
	e := aeridya.Create("./conf")
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	them := &theme{}
	err := them.Init("main.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	aeridya.Theme = them
	e = aeridya.Run()
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

type theme struct {
	statictheme.Theme
	aeridya.Page
}

func (t *theme) Init(s string) error {
	if err := t.StaticInit(s); err != nil {
		return err
	}
	a := template.New("default")
	t.Errors[404], _ = a.Parse("Page not found: {{.Status}}\n{{.Err}}\n")

	t.Pages["/"] = t
	return nil
}

func (t theme) Get(resp *aeridya.Response) {
	resp.Good(200)
	fmt.Fprintf(resp.W, "Hello Test!\n")
	return
}
