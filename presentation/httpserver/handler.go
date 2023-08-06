package httpserver

import "net/http"

type Handler interface {
	GetReq(w http.ResponseWriter, r *http.Request)
}

type handler struct {
}

func NewHttpHandler() Handler {
	return &handler{}
}

func (h *handler) GetReq(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("body yay!")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
