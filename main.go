package main

// @title           Openflat Server API
// @version         0.1
// @description     This is a APIs for OpenflatServer.
// @termsOfService

// @contact.name   Girish Bhutiya
// @contact.url    #
// @contact.email  support@openfiat.org

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/swisscdn/OpenfiatServer/api"
	"github.com/swisscdn/OpenfiatServer/db"
	_ "github.com/swisscdn/OpenfiatServer/docs"
	"github.com/swisscdn/OpenfiatServer/token"
	util "github.com/swisscdn/OpenfiatServer/utils"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("can not create toekn maker: %v", err)
	}

	//session, err := util.GetAstraDBSession(config)
	client, err := util.GetAstraDBClient(config)
	if err != nil {
		log.Fatal("can not connect to DB:", err)
	}

	server := &api.Server{
		Config:     config,
		TokenMaker: tokenMaker,
		Store:      &db.DB{DB: client},
	}
	if config.DBMigrateUp {
		db.MigrateUpDB(&db.DB{DB: client})
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.ListenPort),
		Handler: server.Routes(),
	}

	//start the server
	fmt.Printf("server started on port %s", config.ListenPort)
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
