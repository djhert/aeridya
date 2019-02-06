package router

import (
	"errors"
	"fmt"
	"net/http"
)

type Response struct {
	Status int
	Error  error
	Data   interface{}
}

type Router interface {
	Serve(w http.ResponseWriter, r *http.Request) *Response
	Error(resp *Response, w http.ResponseWriter, r *http.Request) *Response
}

type BasicRoute struct {
}

func (b BasicRoute) Serve(w http.ResponseWriter, r *http.Request) *Response {
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Hello Aeridya!\n")
		return nil
	}
	err := errors.New("BasicRoute only supports /")
	return b.Error(&Response{Status: 404, Error: err, Data: nil}, w, r)
}

func (b BasicRoute) Error(resp *Response, w http.ResponseWriter, r *http.Request) *Response {
	fmt.Fprintf(w, "Error: %d\n%s\n", resp.Status, resp.Error)
	return resp
}
