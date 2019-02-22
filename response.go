package aeridya

import (
	"errors"
	"net/http"
)

type Response struct {
	W      http.ResponseWriter
	R      *http.Request
	Status int
	Err    error
	Data   interface{}
}

func (r *Response) Error(msg string) {
	r.Err = errors.New(msg)
}

func (r *Response) Good(status int) {
	r.Status = status
	r.W.WriteHeader(status)
}

func (r *Response) Bad(status int, msg string) {
	r.Status = status
	r.Error(msg)
	r.W.WriteHeader(status)
}

func (r *Response) Redirect(status int, url string) {
	r.Status = status
	http.Redirect(r.W, r.R, url, status)
}

func (r *Response) GetCookie(cookie string) (*http.Cookie, error) {
	return r.R.Cookie(cookie)
}

func (r *Response) AddCookie(name, value string, maxage int) {
	c := http.Cookie{Name: name, Value: value, MaxAge: maxage}
	http.SetCookie(r.W, &c)
}

func (r *Response) AddRawCookie(c http.Cookie) {
	http.SetCookie(r.W, &c)
}

func (r *Response) DeleteCookie(name string) {
	c := http.Cookie{Name: name, MaxAge: -1}
	http.SetCookie(r.W, &c)
}
