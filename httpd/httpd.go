package httpd

import (
	"github.com/gorilla/mux"
	"github.com/sencoder/go-loghouse/config"
	"github.com/sencoder/go-loghouse/httpd/handler"
	"github.com/sencoder/go-loghouse/httpd/middleware"
	"github.com/sencoder/go-loghouse/pkg/log"
	"net/http"
	"net/http/pprof"
	"time"
)

func newRouter() http.Handler {
	router := mux.NewRouter()
	router.Use(middleware.LogHandler)

	router.Methods("GET").Path("/").HandlerFunc(handler.RootHandler)

	router.Methods("GET").Path("/ping").HandlerFunc(handler.PingHandler)
	router.Methods("GET").Path("/query").HandlerFunc(handler.QueryHandler)
	router.Methods("GET").PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)

	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("build/"))).Methods("GET")
	return router
}

func RunServer(cfg *config.HttpConfig) error {

	r := newRouter()
	server := &http.Server{
		Addr:        cfg.Listen,
		Handler:     r,
		IdleTimeout: time.Second * 10,
	}
	log.Infof("start http server on %s", cfg.Listen)

	return server.ListenAndServe()
}
