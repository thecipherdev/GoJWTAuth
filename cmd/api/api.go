package api

import (
	"net/http"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (r *APIServer) Run() error {
	router := http.NewServeMux()

	return http.ListenAndServe(r.addr, router)
}
