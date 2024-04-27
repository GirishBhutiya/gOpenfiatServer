package routes

import (
	"log"
	"net/http"

	"github.com/GirishBhutiya/gOpenfiatServer/app/handler"
	"github.com/GirishBhutiya/gOpenfiatServer/app/middleware"
	"github.com/GirishBhutiya/gOpenfiatServer/app/token"
	_ "github.com/GirishBhutiya/gOpenfiatServer/docs"

	gochi "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Routes(tokenMaker token.Maker) http.Handler {
	log.Println("Routes")
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

	mux.Use(gochi.Heartbeat("/ping"))

	mux.Post("/", handler.Brocker)
	mux.Get("/", handler.Brocker)

	//mux.Post("/login", app.Login)
	mux.Post("/register", handler.Register)
	mux.Post("/verifyotp", handler.VerifyOTP)

	middleware.InitAuthTokenMaker(&tokenMaker)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(middleware.AuthMiddleware)
		mux.Post("/update", handler.UpdateUser)
		mux.Post("/delete", handler.DeleteUser)
		mux.Post("/create-order", handler.CreateOrder)
		mux.Post("/update-ordervalue", handler.UpdateOrderValue)
		mux.Post("/order-confirming", handler.ConfirmingOrder)
		mux.Post("/order-confirm", handler.ConfirmOrder)
		mux.Post("/order-disputed", handler.DisputedOrder)
		mux.Post("/order-delete", handler.DeleteOrder)
		mux.Post("/allorders", handler.GetUserAllOrders)

	})
	//SwaggerRequest(mux)
	mux.Mount("/swagger", httpSwagger.WrapHandler)
	return mux
}
