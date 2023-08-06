package httpserver

import (
	"context"
	"net/http"
)

type HttpServer interface {
	SetHandleFunc(pattern string, handlerFunc http.HandlerFunc)
	Run() error
	Shutdown(ctx context.Context) error
}

type httpServer struct {
	srv *http.Server
	mux *http.ServeMux
}

func NewHttpServer(port string) HttpServer {
	return &httpServer{
		srv: &http.Server{
			Addr: port,
		},
		mux: http.NewServeMux(),
	}
}

func (hs *httpServer) SetHandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	// register to mux so and don't use http.HandleFunc for multiple registration (testing)
	hs.mux.HandleFunc(pattern, handlerFunc)
}

func (hs *httpServer) Run() error {
	hs.srv.Handler = hs.mux
	return hs.srv.ListenAndServe()
}

func (hs *httpServer) Shutdown(ctx context.Context) error {
	return hs.srv.Shutdown(ctx)
}
