package main

import (
	"net/http"
	"os"
	"time"

	"github.com/diontristen/go-crud/api"
	"github.com/diontristen/go-crud/util"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/muravjov/slog/watcher"
)

func addContext(f func(w http.ResponseWriter, r *http.Request, ac *util.AppContext), ac *util.AppContext) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r, ac)
	}
}

func appMain() (exitOk bool) {
	config, err := util.NewConfig()
	if err != nil {
		util.Error(err)
		return
	}

	ac, acError := util.NewAppContext(config)

	if acError != nil {
		util.Errorf("error when creating context: %v\n", acError)
		return
	}

	defer ac.Close()

	router := mux.NewRouter()
	router.HandleFunc("/v1/contact", addContext(api.GetContact, ac)).Methods(http.MethodGet)
	router.HandleFunc("/v1/contact", addContext(api.PostContact, ac)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/v1/contact", addContext(api.RemoveContact, ac)).Methods(http.MethodDelete)
	router.HandleFunc("/v1/update/contact", addContext(api.UpdateContact, ac)).Methods(http.MethodPost, http.MethodOptions)

	http.Handle("/", router)

	server := &http.Server{Addr: getListenAddr(ac.Config), Handler: &util.ServerHandler{
		Router: router,
	}}
	util.Infof("App listening address %s", server.Addr)
	return util.ListenAndServe(server, func() {})
}

func getListenAddr(config util.Config) string {
	listen := ":5000"
	if configListen := config.Listen; configListen != nil {
		listen = *configListen
	}
	return listen
}

func setupLogging(config util.Config) func() {
	close := func() {}

	if config.DebugFlag {
		util.DebugFlag = true
	}

	if config.SentryDSN != "" {
		watcher.StartWatcher(config.SentryDSN, "")

		err := util.SetupSentry(config.SentryDSN)
		if err != nil {
			util.Errorf("SetupSentry: %s", err)
		} else {
			close = func() {
				if !sentry.Flush(time.Second * 10) {
					util.Error("Some Sentry events may not have been sent")
				}
			}
		}
	}

	return close
}

func main() {
	code := 0
	if !appMain() {
		code = 1
	}

	os.Exit(code)
}
