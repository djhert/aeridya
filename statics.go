// Allows static items to be served from domain.com/item
// Items requested should be located at root of statics dir
package aeridya

import (
	"github.com/hlfstr/aeridya/handler"
	"net/http"
)

// Statics is an object that embeds Aeridya's Handler into it to allow for a list of
// handler functions to process each static file.
// Dir is the directory where statics are on the file system
// Statics.Statics is a slice of strings to serve via this system
type statics struct {
	handler.Handler
	Dir     string
	Statics []string
}

// Create will create the Statics Object and return it
// Requires a path to the statics dir
func NewStatics(path string) *statics {
	s := &statics{handler.Handler{}, path, make([]string, 0)}
	s.Init()
	return s
}

// Defaults places the default items all domains need into statics array
// Meaning:  favicon.ico, sitemap.xml, and robots.txt
func (s *statics) Defaults() {
	s.Statics = append(s.Statics, "/favicon.ico")
	s.Statics = append(s.Statics, "/sitemap.xml")
	s.Statics = append(s.Statics, "/robots.txt")
}

// Add requires an item, adds item into Statics array
func (s *statics) Add(item string) {
	s.Statics = append(s.Statics, item)
}

// serve creates a http.Handler and adds it to the DefaultServeMux of the application
func (s statics) serve(pattern string, filename string) {
	http.Handle(pattern, s.Final(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})))
}

// Serve iterates over Statics array and creates a handler for each static object
// Requires array of http.Handlers for Aeridya Handler compatibility
func (s statics) Serve(handles []func(http.Handler) http.Handler) {
	s.AppendHandlers(handles)
	for t := range s.Statics {
		s.serve(s.Statics[t], s.Dir+s.Statics[t])
	}
}
