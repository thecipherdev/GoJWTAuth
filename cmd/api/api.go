package api

import (
	"log"
	"net/http"

	"github.com/thecipherdev/goauth/controller"
	"github.com/thecipherdev/goauth/middleware"
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
	mainRouter := http.NewServeMux()

	userRouter := http.NewServeMux()
	userHandler := controller.NewUserHandler()
	userHandler.UserRouter(userRouter)

	mainRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", userRouter))

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.IsAuthenticated,
	)

	server := http.Server{
		Addr:    r.addr,
		Handler: stack(mainRouter),
	}

	log.Printf("Server is running on PORT: %v", r.addr)
	return server.ListenAndServe()
}
