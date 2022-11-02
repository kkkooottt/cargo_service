package main

import (
	"cargo_service/server"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	mySecretStringForJWT = "strongSecretString"
	serverPort           = "8080"
	DBhost               = "localhost"
	DBport               = 5432
	DBuser               = "kkkooottt"
	DBpassword           = "secretpassword"
)

/*
// token(without expiring) eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJyb2xlIjoiYWRtaW4iLCJ1c2VyIjoidXNlcm5hbWUifQ.bNdNW0Jd2QkrQIcwwPoFsmk2HpBlroJzTjb1PGy1I6s
{
  "authorized": true,
  "role": "admin",
  "user": "username"
}
*/

// create // curl -X POST -i -d '{"user":"","weight":45,"volume":67,"load_date_time":"0001-01-01T00:00:00Z","unload_date_time":"0001-01-01T00:00:00Z","load_adress":"moscow","unload_adress":"moon","comment":"tratata"}' -H 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJyb2xlIjoiYWRtaW4iLCJ1c2VyIjoidXNlcm5hbWUifQ.bNdNW0Jd2QkrQIcwwPoFsmk2HpBlroJzTjb1PGy1I6s' 0.0.0.0:8080/deals/create
// id info // curl -i -H 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJyb2xlIjoiYWRtaW4iLCJ1c2VyIjoidXNlcm5hbWUifQ.bNdNW0Jd2QkrQIcwwPoFsmk2HpBlroJzTjb1PGy1I6s' 0.0.0.0:8080/deals/1/info
// list  //  curl -i -H 'Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJyb2xlIjoiYWRtaW4iLCJ1c2VyIjoidXNlcm5hbWUifQ.bNdNW0Jd2QkrQIcwwPoFsmk2HpBlroJzTjb1PGy1I6s' "0.0.0.0:8080/deals/list?page=1&perpage=1"

func main() {

	db := initDB()
	fmt.Println("ping db:", db.Ping())

	server.StartServer(db, serverPort, mySecretStringForJWT)

}

func initDB() *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s sslmode=disable",
		DBhost, DBport, DBuser, DBpassword)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE test")
	if err != nil {
		fmt.Println(err)
	}
	//(вес груза, объем груза, дата+время погрузки, дата+время выгрузки, адрес погрузки, адрес выгрузки, комментарий к заказу)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS deals (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255),
		weight INT,
		volume INT,
		load_datetime timestamp,
		unload_datetime timestamp,
		load_adress TEXT,
		unload_adress TEXT,
		comment TEXT
	);`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db.Ping())
	return db

}
