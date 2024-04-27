package api

import (
	"log"
	"net/http"

	"github.com/GirishBhutiya/gOpenfiatServer/app/config"
	db "github.com/GirishBhutiya/gOpenfiatServer/app/database"
	"github.com/GirishBhutiya/gOpenfiatServer/app/handler"
	"github.com/GirishBhutiya/gOpenfiatServer/app/routes"
	"github.com/GirishBhutiya/gOpenfiatServer/app/token"
)

var (
	server *handler.Server
)

// @title           Openflat Server API
// @version         0.1
// @description     This is a APIs for OpenflatServer.
// @termsOfService

// @contact.name   openfiat.org
// @contact.url    #
// @contact.email  support@openfiat.org

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      g-openfiat-server.vercel.app
// @BasePath  /
func init() {
	con, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	tokenMaker, err := token.NewPasetoMaker(con.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("can not create toekn maker: %v", err)
	}

	//session, err := util.GetAstraDBSession(config)
	client, err := config.GetAstraDBClient(con)
	if err != nil {
		log.Fatal("can not connect to DB:", err)
	}

	server = &handler.Server{
		Config:     con,
		TokenMaker: tokenMaker,
		Store:      &db.DB{DB: client},
	}
	if con.DBMigrateUp {
		db.MigrateUpDB(&db.DB{DB: client})
	}
	handler.InitServer(server)
	log.Println("Vercel Init complete")

}

// Entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Vercel handler")
	//app.ServeHTTP(w, r)
	routes.Routes(server.TokenMaker).ServeHTTP(w, r)
}
