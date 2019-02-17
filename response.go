package aeridya

import (
	"errors"
	"net/http"
)

type Response struct {
	Status int
	Err    error
	Data   interface{}
}

func (r *Response) Error(msg string) {
	r.Err = errors.New(msg)
}

func (r *Response) Good(status int, w http.ResponseWriter) {
	r.Status = status
	w.WriteHeader(status)
}

func (r *Response) Bad(status int, msg string, w http.ResponseWriter) {
	r.Status = status
	r.Error(msg)
	w.WriteHeader(status)
}

func (s *Response) Redirect(status int, url string, w http.ResponseWriter, r *http.Request) {
	s.Status = status
	http.Redirect(w, r, url, status)
}
