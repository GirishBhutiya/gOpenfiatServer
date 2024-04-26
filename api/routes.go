package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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
		mux.Post("/delete", app.DeleteUser)
		mux.Post("/create-order", app.CreateOrder)
		mux.Post("/update-ordervalue", app.UpdateOrderValue)
		mux.Post("/order-confirming", app.ConfirmingOrder)
		mux.Post("/order-confirm", app.ConfirmOrder)
		mux.Post("/order-disputed", app.DisputedOrder)
		mux.Post("/order-delete", app.DeleteOrder)
		mux.Post("/allorders", app.GetUserAllOrders)

	})
	//SwaggerRequest(mux)
	mux.Mount("/swagger", httpSwagger.WrapHandler)
	return mux
}
