package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/swisscdn/OpenfiatServer/db"
	"github.com/swisscdn/OpenfiatServer/token"
	util "github.com/swisscdn/OpenfiatServer/utils"
)

// Handler is the entrypoint for the vercel serverless function
func Handler(w http.ResponseWriter, req *http.Request) {
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

	server := Server{
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
