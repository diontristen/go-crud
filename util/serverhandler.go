package util

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
)

type ServerHandler struct {
	Router *mux.Router
}

type Exception struct {
	Type  string      `json:"type"`
	Error interface{} `json:"error"`
}

// USAGE: Returns 404 Response code if route doesn't match route defined against the route. (application.go/appMain())
func (handler *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// do not use panic(nil)
		if err := recover(); err != nil {
			// copy-n-paste from http/net/server.go
			// for more clear file logging
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			Errorf("ServeHTTP exception: %s\n%s", err, buf)

			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(Exception{
				Type:  "exception",
				Error: err,
			}); err != nil {
				Errorf("json.Encode: %s", err)
			}
		}
	}()

	var match mux.RouteMatch
	if handler.Router.Match(r, &match) {
		r = mux.SetURLVars(r, match.Vars)
		match.Handler.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

// USAGE: Starts the http server with the given server and handle (mux Router)
// Checks if a SIGINT or SIGTERM is triggered in the command line
func ListenAndServe(server *http.Server, beforeShutdown func()) bool {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	var servicesWg sync.WaitGroup
	serverOk := true

	servicesWg.Add(1)
	go func() {
		defer servicesWg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Error(err)
			serverOk = false
		}
	}()

	Info("waiting for termination signal...")
	sig := <-sigChan
	Infof("signal received: %s", sig.String())

	beforeShutdown()

	if err := server.Shutdown(context.Background()); err != nil {
		Errorf("HTTP server Shutdown: %v", err)
		return false
	}

	servicesWg.Wait()
	return serverOk
}
