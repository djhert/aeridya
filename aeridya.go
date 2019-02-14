package aeridya

import (
	"fmt"
	"github.com/hlfstr/aeridya/handler"
	"github.com/hlfstr/aeridya/router"
	"github.com/hlfstr/aeridya/statics"
	"github.com/hlfstr/configurit"
	"github.com/hlfstr/logit"
	"net/http"
	"os"
	"runtime"
)

//Global reference to Aeridya instance
var instance *Aeridya

//Create the instance of Aeridya at init
func init() {
	instance = new(Aeridya)
}

//Aeridya Type Definition
type Aeridya struct {
	Logger  *logit.Logger
	Handler *handler.Handler
	Statics *statics.Statics

	BaseTemplate string
	TemplateDir  string
	Port         string
	Domain       string
	Development  bool

	Route router.Router
}

func Create(conf string) (*Aeridya, *configurit.Conf, error) {
	instance = new(Aeridya)
	c, err := instance.loadConfig(conf)
	if err != nil {
		return nil, nil, err
	}
	instance.Statics.Defaults()
	instance.Handler = handler.Create()
	instance.Route = router.BasicRoute{}
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
			if a.Logger, err = logit.StartLogger(logit.TermLog()); err != nil {
				return nil, err
			}
		} else {
			if file, err := logit.OpenFile(log); err != nil {
				return nil, err
			} else {
				if a.Logger, err = logit.StartLogger(file); err != nil {
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

	if t, err := c.GetString("", "Templates"); err != nil {
		return nil, err
	} else {
		a.TemplateDir = t
	}

	if a.Development, err = c.GetBool("", "Development"); err != nil {
		return nil, err
	}

	return c, err
}

func (a *Aeridya) Run() error {
	defer a.panicAttack()
	a.Logger.Logf("Starting %s for %s | Listening on %s", Version(), a.Domain, a.Port)
	http.Handle("/", a.Handler.Final(http.HandlerFunc(a.ServeHTTP)))
	go a.Statics.Serve(a.Handler.Get())
	return http.ListenAndServe(a.Port, nil)
}

func (a Aeridya) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if resp := a.Route.Serve(w, r); resp.Error != nil {
		if a.Development {
			a.Logger.Logf("[Error[%d]] %s", resp.Status, resp.Error.Error())
		}
	}
	return
}

func mkerror(msg string) error {
	return fmt.Errorf("Aeridya: %s", msg)
}

func (a Aeridya) panicAttack() {
	err := recover()
	if err != nil {
		a.Logger.Logf("PANIC!\n  %#v\n", err)
		buf := make([]byte, 4096)
		buf = buf[:runtime.Stack(buf, true)]
		a.Logger.Logf("Stack Trace:\n%s\n", buf)
		os.Exit(1)
	}
}
