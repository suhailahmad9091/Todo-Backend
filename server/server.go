package server

import (
	"Todo/handlers"
	"Todo/middlewares"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Use(middlewares.CommonMiddlewares()...)

	router.Route("/v1", func(v1 chi.Router) {
		v1.Post("/register", handlers.RegisterUser)
		v1.Post("/login", handlers.LoginUser)

		v1.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate)

			r.Route("/user", func(user chi.Router) {
				user.Get("/profile", handlers.GetUser)
				user.Post("/logout", handlers.LogoutUser)
				user.Delete("/delete", handlers.DeleteUser)
			})

			r.Route("/todo", func(todo chi.Router) {
				todo.Post("/", handlers.CreateTodo)
				todo.Get("/", handlers.GetAllTodos)
				todo.Delete("/delete-all", handlers.DeleteAllTodos)

				todo.Route("/{todoId}", func(todoIDRoute chi.Router) {
					todoIDRoute.Delete("/", handlers.DeleteTodo)
					todoIDRoute.Put("/mark-completed", handlers.MarkCompleted)
				})
			})
		})
	})

	return &Server{
		Router: router,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
