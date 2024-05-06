package routes

import (
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
	middleware.InitAuthTokenMaker(&tokenMaker)
	mux.Use(gochi.Heartbeat("/ping"))

	mux.Post("/", handler.Brocker)
	mux.Get("/", handler.Brocker)

	fileServer := http.FileServer(http.Dir("./profilepic/"))
	mux.Handle("/profilepic/*", http.StripPrefix("/profilepic", fileServer))

	//mux.Post("/login", app.Login)
	mux.Post("/register", handler.Register)
	mux.Post("/login", handler.Login)
	mux.Post("/renew-accesstoken", handler.RenewAccessToken)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(middleware.AuthMiddleware)
		//user endpoints
		mux.Post("/update", handler.UpdateUser)
		mux.Post("/subscribe", handler.SubscribeGroupToUSer)
		mux.Post("/unsubscribe", handler.UnsubscribeGroupToUSer)
		mux.Post("/getgroups", handler.GetAllGroups)

		//group endpoints
		mux.Post("/create-group", handler.CreateGroup)
		mux.Post("/update-group", handler.UpdateGroup)
		mux.Post("/delete-group", handler.DeleteGroup)

		//order endpoints
		mux.Post("/create-buy-order", handler.CreateBuyOrder)
		mux.Post("/create-sell-order", handler.CreateSellOrder)
		mux.Post("/update-order", handler.UpdateOrder)
		mux.Post("/order-delete", handler.DeleteOrder)
		mux.Post("/getorders", handler.GetAllOrders)

		//trade endpoints
		mux.Post("/create-trade", handler.CreateTrade)
		mux.Post("/trade-confirming", handler.ConfirmingTrade)
		mux.Post("/trade-confirm", handler.ConfirmTrade)
		mux.Post("/trade-disputed", handler.DisputedTrade)
		mux.Post("/trade-delete", handler.DeleteTrade)
		mux.Post("/getusertrades", handler.GetAllUsersTrade)
		mux.Post("/getordertrades", handler.GetOrderTrades)
		//mux.Post("/allorders", handler.GetUserAllOrders)

	})
	//SwaggerRequest(mux)
	mux.Mount("/swagger", httpSwagger.WrapHandler)
	return mux
}
