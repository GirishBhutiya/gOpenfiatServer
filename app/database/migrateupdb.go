package db

import (
	"fmt"
	"log"
)

func MigrateUpDB(db *DB) {

	var dropUserDB = db.DB.Query(`DROP TABLE IF EXISTS users`)
	var dropOrdersDB = db.DB.Query(`DROP TABLE IF EXISTS orders`)
	var dropGroupDB = db.DB.Query(`DROP TABLE IF EXISTS group`)
	var dropTradeDB = db.DB.Query(`DROP TABLE IF EXISTS trade`)
	var createUserDB = db.DB.Query(`CREATE TABLE IF NOT EXISTS users ( userid uuid PRIMARY KEY,phonenumber bigint, verified boolean, first_name text,last_name text, groups map<uuid, text>, profile_pic text, otp int);`)
	var createOrderDB = db.DB.Query(`CREATE TABLE IF NOT EXISTS orders (orderid uuid PRIMARY KEY, userid uuid, fiatAmount float, minAmount float, price float,timeLimit timestamp, type text);`)
	var createGroupDB = db.DB.Query(`CREATE TABLE IF NOT EXISTS group( groupid uuid PRIMARY KEY, groupname text, creatorUserid uuid,createtime timestamp);`)
	var createTradeDB = db.DB.Query(`CREATE TABLE IF NOT EXISTS trade( tradeid uuid PRIMARY KEY, orderid uuid, bidUserid uuid,tradetime timestamp,status text, method text);`)

	fmt.Println("Dropping the tables now")
	_, err := dropUserDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = dropOrdersDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = dropGroupDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = dropTradeDB.Exec()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating the tables now")
	_, err = createUserDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = createOrderDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = createGroupDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = createTradeDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
