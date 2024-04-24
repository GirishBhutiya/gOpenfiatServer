package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *Server) Routes() http.Handler {
	mux := chi.NewRouter()

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Brocker)
	mux.Get("/", app.Brocker)

	//mux.Post("/login", app.Login)
	mux.Post("/register", app.Register)
	mux.Post("/verifyotp", app.VerifyOTP)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.AuthMiddleware)
		mux.Post("/update", app.UpdateUser)
		mux.Post("/delete-account", app.DeleteUser)

	})
	//SwaggerRequest(mux)
	mux.Mount("/swagger", httpSwagger.WrapHandler)
	return mux
}
