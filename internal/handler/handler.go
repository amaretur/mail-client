package handler

import (
	"github.com/gorilla/mux"
)

type handler interface {
	Init(router *mux.Router)
}

type Handler struct {
	router *mux.Router
}

func NewHandler() *Handler {

	r := mux.NewRouter().StrictSlash(true).PathPrefix("/api/v1").Subrouter()

	return &Handler{
		router : r,
	}
}

func (h *Handler) Router() *mux.Router {
	return h.router
}

func (h *Handler) Register(sh handler, prefix string) { 
	sh.Init(h.router.PathPrefix(prefix).Subrouter())
}
