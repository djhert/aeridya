package aeridya

import (
	"bytes"
	"html/template"
	"net/http"
)

type errors struct {
	pages map[int]*template.Template
}

type ErrData struct {
	Status int
	Error  error
}

func CreateError() *errors {
	e := &errors{}
	e.pages = make(map[int]*template.Template)

	t := template.New("default")
	e.pages[0], _ = t.Parse("An error has occurred: {{.Status}}\n{{.Error}}")

	if instance.Development {
		t := template.New("log")
		e.pages[1], _ = t.Parse("Error:[{{.Status}}] {{.Error}}")
	}

	tea := template.New("teapot")
	e.pages[418], _ = tea.Parse("Error: {{.Status}}\nI am a teapot!\n")
	return e
}

func (e errors) Show(status int, data interface{}, w http.ResponseWriter, r *http.Request) {
	if t, ok := e.pages[status]; ok {
		t.Execute(w, data)
	} else {
		e.pages[0].Execute(w, data)
	}
	if instance.Development {
		var out bytes.Buffer
		e.pages[1].Execute(&out, data)
		instance.Logger.Log(out.String())
	}
	return
}
