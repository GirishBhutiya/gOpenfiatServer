package db

import (
	"fmt"
	"log"
)

func MigrateUpDB(db *DB) {

	var dropUserDB = db.DB.Query(`DROP TABLE IF EXISTS users`)
	var dropRollDB = db.DB.Query(`DROP TABLE IF EXISTS roll`)
	var createUserDB = db.DB.Query(`CREATE TABLE IF NOT EXISTS users ( phonenumber bigint PRIMARY KEY, verified boolean, first_name text,last_name text, roll_id int, profile_pic text, otp int);`)
	var createRollDB = db.DB.Query(`CREATE TABLE IF NOT EXISTS roll ( id int PRIMARY KEY, roll text);`)
	var addAdminRoll = db.DB.Query(`INSERT INTO roll(id,roll) VALUES (1,'admin');`)
	var addUserRoll = db.DB.Query(`INSERT INTO roll(id,roll) VALUES (2,'user');`)

	fmt.Println("Dropping the tables now")
	_, err := dropUserDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = dropRollDB.Exec()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Creating the tables now")
	_, err = createRollDB.Exec()
	if err != nil {
		log.Fatal(err)
	}
	_, err = createUserDB.Exec()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Insering the data into tables now")
	err = db.DB.Batch(addAdminRoll, addUserRoll).Exec()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.DB.Query("select * FROM roll;").Exec()
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range rows {
		vals := r.Values()
		fmt.Println("Values:", vals)
	}
}
