package http

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	router       http.Handler
	port         int
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewServer(port int, readTimeout, writeTimeout time.Duration, router http.Handler) (*Server, error) {
	return &Server{
		router:       router,
		port:         port,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}, nil
}

func (a *Server) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *Server) Run() error {

	s := &http.Server{
		Addr:         fmt.Sprintf(":%v", a.port),
		Handler:      a.router,
		ReadTimeout:  a.readTimeout,
		WriteTimeout: a.writeTimeout,
	}

	return s.ListenAndServe()
}
