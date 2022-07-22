package server

import (
	"ToDo/handlers"
	"ToDo/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	chi.Router
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Post("/signup", handlers.Signup)
		api.Post("/login", handlers.Login)

		api.Route("/todo", func(r chi.Router) {
			r.Use(middleware.Middleware)
			r.Use(middleware.Recovery)
			r.Get("/logout", handlers.SignOut)
			r.Post("/task", handlers.AddTask)
			r.Get("/all-task", handlers.FetchAllTask)
			r.Route("/{id}", func(ap chi.Router) {
				ap.Delete("/task", handlers.DeleteTask)
				ap.Put("/", handlers.UpdateTask)
				ap.Put("/complete", handlers.CompleteTask)
			})

		})
	})

	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
