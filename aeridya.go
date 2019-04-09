package aeridya

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/securecookie"
	"github.com/hlfstr/aeridya/handler"
	"github.com/hlfstr/configurit"
	"github.com/hlfstr/logit"
	"github.com/hlfstr/quit"
	"github.com/hlfstr/quit/quitters"
)

var (
	server *http.Server

	Log     *logit.Logger
	Handler *handler.Handler
	Static  *statics
	Config  *configurit.Conf

	Port        string
	Domain      string
	FullDomain  string
	Development bool
	HTTPS       bool

	Theme theming

	cookieHash  []byte
	cookieBlock []byte

	isInit  bool
	limiter chan struct{}
)

func Create(conf string) error {
	var err error
	Config, err = loadConfig(conf)
	if err != nil {
		return err
	}
	//quitter.AddQuit(shutdown)
	quitters.AddQuit(quit.Quit)
	quitters.AddQuit(Log.Quit)
	Static.Defaults()
	Handler = handler.Create()
	Theme = &ATheme{}
	cookieHandler = securecookie.New(cookieHash, cookieBlock)
	if HTTPS {
		FullDomain = "https://" + Domain
	} else {
		FullDomain = "http://" + Domain
	}
	server = &http.Server{Addr: Port}
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

	if h, err := c.GetString("", "CookieHash"); err != nil {
		return nil, err
	} else {
		cookieHash = []byte(h)
	}

	if h, err := c.GetString("", "CookieBlock"); err != nil {
		return nil, err
	} else {
		cookieBlock = []byte(h)
	}

	if log, err := c.GetString("", "Log"); err != nil {
		return nil, err
	} else {
		if log == "stdout" {
			if Log, err = logit.Start(logit.TermLog()); err != nil {
				return nil, err
			}
		} else {
			if file, err := logit.OpenFile(log); err != nil {
				return nil, err
			} else {
				if Log, err = logit.Start(file); err != nil {
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

	if HTTPS, err = c.GetBool("", "HTTPS"); err != nil {
		return nil, err
	}

	return c, err
}

func Run() error {
	if !isInit {
		return mkerror("Must use Create(\"/path/to/config\") before Run()")
	}
	go quit.Run(shutdown, shutdown)
	defer panicAttack()
	defer quitters.Quit()
	Log.Logf(logit.MSG, "Starting %s for %s | Listening on %s", Version(), Domain, Port)
	http.Handle("/", Handler.Final(internalTrailingSlash(http.HandlerFunc(serve))))
	go Static.Serve(Handler.Get())
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		Log.LogError(1, err)
		return nil
	}
	Log.Log(0, "Closing aeridya")
	return nil
}

func serve(w http.ResponseWriter, r *http.Request) {
	resp := &Response{W: w, R: r}
	Theme.Serve(resp)
	if resp.Err != nil {
		if Development {
			Log.Logf(logit.ERROR, "[Error(%d)] %s", resp.Status, resp.Err.Error())
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

func shutdown() {
	if err := server.Shutdown(context.Background()); err != nil {
		Log.LogError(1, err)
		server.Close()
	}
}

func panicAttack() {
	err := recover()
	if err != nil {
		Log.Logf(logit.PANIC, "PANIC!\n  %#v\n", err)
		buf := make([]byte, 4096)
		buf = buf[:runtime.Stack(buf, true)]
		Log.Logf(logit.PANIC, "Stack Trace:\n%s\n", buf)
		os.Exit(1)
	}
}
