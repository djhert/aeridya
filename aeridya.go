package aeridya

import (
	"fmt"
	"github.com/hlfstr/aeridya/handler"
	"github.com/hlfstr/configurit"
	"github.com/hlfstr/logit"
	"net/http"
	"os"
	"runtime"
)

var (
	Logger  *logit.Logger
	Handler *handler.Handler
	Static  *statics
	Config  *configurit.Conf

	Port        string
	Domain      string
	Development bool

	Theme theming

	quitters []func()

	isInit  bool
	limiter chan struct{}
)

func Create(conf string) error {
	var err error
	Config, err = loadConfig(conf)
	if err != nil {
		return err
	}
	quitters = make([]func(), 0)
	AddQuit(Logger.Quit)
	Static.Defaults()
	Handler = handler.Create()
	Theme = &ATheme{}
	isInit = true
	return nil
}

func loadConfig(conf string) (*configurit.Conf, error) {
	c, err := configurit.Open(conf)
	if err != nil {
		return nil, err
	}

	if Domain, err = c.GetString("", "Domain"); err != nil {
		return nil, err
	}

	if s, err := c.GetString("", "Port"); err != nil {
		return nil, err
	} else {
		Port = ":" + s
	}

	if n, err := c.GetInt("", "Workers"); err != nil {
		return nil, err
	} else {
		limiter = make(chan struct{}, n)
	}

	if log, err := c.GetString("", "Log"); err != nil {
		return nil, err
	} else {
		if log == "stdout" {
			if Logger, err = logit.StartLogger(logit.TermLog()); err != nil {
				return nil, err
			}
		} else {
			if file, err := logit.OpenFile(log); err != nil {
				return nil, err
			} else {
				if Logger, err = logit.StartLogger(file); err != nil {
					return nil, err
				}
			}
		}
	}
	if s, err := c.GetString("", "Statics"); err != nil {
		return nil, err
	} else {
		Static = NewStatics(s)
	}
	if Development, err = c.GetBool("", "Development"); err != nil {
		return nil, err
	}

	return c, err
}

func AddQuit(f func()) {
	quitters = append(quitters, f)
}

func Run() error {
	if !isInit {
		return mkerror("Must use Create(\"/path/to/config\") before Run()")
	}
	defer panicAttack()
	defer quit()
	Logger.Logf("Starting %s for %s | Listening on %s", Version(), Domain, Port)
	http.Handle("/", Handler.Final(internalTrailingSlash(http.HandlerFunc(serve))))
	go Static.Serve(Handler.Get())
	return http.ListenAndServe(Port, nil)
}

func serve(w http.ResponseWriter, r *http.Request) {
	resp := &Response{W: w, R: r}
	Theme.Serve(resp)
	if resp.Err != nil {
		if Development {
			Logger.Logf("[Error(%d)] %s", resp.Status, resp.Err.Error())
		}
	}
	return
}

func limit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter <- struct{}{}
		h.ServeHTTP(w, r)
		<-limiter
	})
}

func AddTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Path
		if u[len(u)-1:] != "/" {
			u = u + "/"
			http.Redirect(w, r, u, 301)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func NoTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Path
		if u[len(u)-1:] == "/" {
			u = u[:len(u)-1]
			http.Redirect(w, r, u, 301)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func internalTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Path
		if len(u) > 1 {
			if u[len(u)-1:] == "/" {
				u = u[:len(u)-1]
			}
		} else {
			u = "/"
		}
		r.URL.Path = u
		h.ServeHTTP(w, r)
	})
}

func mkerror(msg string) error {
	return fmt.Errorf("Aeridya[Error]: %s", msg)
}

func quit() {
	for i := range quitters {
		quitters[i]()
	}
}

func panicAttack() {
	err := recover()
	if err != nil {
		Logger.Logf("PANIC!\n  %#v\n", err)
		buf := make([]byte, 4096)
		buf = buf[:runtime.Stack(buf, true)]
		Logger.Logf("Stack Trace:\n%s\n", buf)
		os.Exit(1)
	}
}
