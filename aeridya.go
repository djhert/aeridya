package aeridya

import (
	"fmt"
	"github.com/hlfstr/aeridya/handler"
	//	"github.com/hlfstr/aeridya/router"
	"github.com/hlfstr/aeridya/statics"
	"github.com/hlfstr/configurit"
	"github.com/hlfstr/logit"
	"net/http"
)

//Global reference to Aeridya instance
var instance *Aeridya

//Create the instance of Aeridya at init
func init() {
	instance = new(Aeridya)
}

//Aeridya Type Definition
type Aeridya struct {
	Logger         *Logit.Logger
	DefaultHandler *handler.Handler
	Statics        *statics.Statics

	Port        string
	Domain      string
	Development bool

	//	Route *route.Handler
	Errors *errors
}

func Create(conf string) (*Aeridya, *configurit.Conf, error) {
	if instance == nil {
		instance = new(Aeridya)
	}
	c, err := instance.loadConfig(conf)
	if err != nil {
		return nil, nil, err
	}
	instance.Statics.Defaults()
	instance.Errors = CreateError()
	instance.DefaultHandler = handler.Create()
	return instance, c, nil
}

func (a *Aeridya) loadConfig(conf string) (*configurit.Conf, error) {
	c, err := configurit.Open(conf)
	if err != nil {
		return nil, err
	}

	if a.Domain, err = c.GetString("", "Domain"); err != nil {
		return nil, err
	}

	if s, err := c.GetString("", "Port"); err != nil {
		return nil, err
	} else {
		a.Port = ":" + s
	}

	if log, err := c.GetString("", "Log"); err != nil {
		return nil, err
	} else {
		if log == "stdout" {
			if a.Logger, err = Logit.StartLogger(Logit.TermLog()); err != nil {
				return nil, err
			}
		} else {
			if file, err := Logit.OpenFile(log); err != nil {
				return nil, err
			} else {
				if a.Logger, err = Logit.StartLogger(file); err != nil {
					return nil, err
				}
			}
		}
	}

	if sdir, err := c.GetString("", "Statics"); err != nil {
		return nil, err
	} else {
		a.Statics = statics.Create(sdir)
	}

	if a.Development, err = c.GetBool("", "Development"); err != nil {
		return nil, err
	}

	return c, err
}

func (a *Aeridya) Run() error {
	a.Logger.Logf("Starting %s for %s | Listening on %s", Version(), a.Domain, a.Port)
	http.Handle("/", a)
	a.Statics.Serve(a.DefaultHandler.Get())
	return http.ListenAndServe(a.Port, nil)
}

func (a Aeridya) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*if len(r.URL.Path) > 1 {
		temp := r.URL.Path
		if temp[len(temp)-1:] == "/" {
			temp = temp[1 : len(temp)-1]
		} else {
			temp = temp[1:]
		}
		parsed := strings.Split(temp, "/")
		if val, ok := route.pages[parsed[0]]; ok {
			val.Run(w, r, parsed)
			return
		} else {
			if route.customRoute(parsed, w, r) {
				return
			}
			route.Errors.showError(404, w, r)
			return

		}
	} else {
		route.pages["/"].Run(w, r, nil)
		return
	}*/
	a.Errors.Show(404, nil, w, r)
	return
}

func mkerror(msg string) error {
	return fmt.Errorf("Aeridya: %s", msg)
}
