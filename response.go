package aeridya

import (
	"errors"
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

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

func (r *Response) GetCookieValues(name string) (map[string]string, error) {
	ck, err := r.R.Cookie(name)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	err = cookieHandler.Decode(name, ck.Value, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (r *Response) AddCookie(name string, hour int, values map[string]string) error {
	enc, err := cookieHandler.Encode(name, values)
	if err != nil {
		return err
	}
	c := http.Cookie{
		Name:     name,
		Value:    enc,
		Path:     "/",
		MaxAge:   60 * hour,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(r.W, &c)
	return nil
}

func (r *Response) DeleteCookie(name string) {
	c := http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(r.W, &c)
}
