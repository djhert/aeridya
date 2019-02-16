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

	isInit  bool
	limiter chan struct{}
)

func Create(conf string) error {
	var err error
	Config, err = loadConfig(conf)
	if err != nil {
		return err
	}
	Static.Defaults()
	Handler = handler.Create()
	//	Theme = &Theming{}
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

func Run() error {
	if !isInit {
		return mkerror("Must use Create(\"/path/to/config\") before Run()")
	}
	defer panicAttack()
	Logger.Logf("Starting %s for %s | Listening on %s", Version(), Domain, Port)
	http.Handle("/", Handler.Final(limit(http.HandlerFunc(serve))))
	go Static.Serve(Handler.Get())
	return http.ListenAndServe(Port, nil)
}

func serve(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	Theme.Serve(w, r, resp)
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
		defer func() { <-limiter }()
		h.ServeHTTP(w, r)
	})
}

func mkerror(msg string) error {
	return fmt.Errorf("Aeridya[Error]: %s", msg)
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
