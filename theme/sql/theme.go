package sqltheme

import (
	"errors"
	"github.com/hlfstr/aeridya"
	"database/sql"
	"html/template"
)

type Theme struct {
	DB *sql.DB
}

func (t *Theme)SQLInit() error {
	return errors.New("SQL Theme Requires a custom function for SQLInit to connect to the Database.")
}

func (t Theme) Parse(input string, resp *aeridya.Response) *template.HTML {
	out := "hello"
	return &(template.HTML)out
}

func (t Theme) Serve(resp *aeridya.Response) {
	o := t.Parse(resp.R.URL.Path, resp)
	if o == nil {
		resp.Bad(404, resp.Err.Error())
		aeridya.ThemeError(resp, t)
		return
	}
	aeridya.ServePage(resp, o)
}

func (t Theme) Error(resp *aeridya.Response) {

}