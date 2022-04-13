package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	DB_HOST     = "database" //host.docker.internal
	DB_PORT     = 5432
	DB_USER     = "base_null"
	DB_PASSWORD = "base_null"
	DB_NAME     = "base_null"
)

type Ways struct {
	id       int
	char_id  string `json:"char_id"`
	location string `json:"location"`
	the_time string `json:"time"`
}

type JSONresponse struct {
	Type    string `json:"type"`
	Data    []Ways `json:"data"`
	Message string `json:"message"`
}

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func setupDB() *sql.DB {

	DB_USER = GetEnv("DB_USER", "base_null")
	DB_PASSWORD = GetEnv("DB_PASSWORD", "base_null")
	DB_NAME = GetEnv("DB_NAME", "base_null")

	if DB_USER == "base_null" || DB_PASSWORD == "base_null" || DB_NAME == "base_null" {
		panic("One of DB enviroment not found")
	}

	//connStr := "user=" + DB_USER + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " sslmode=disable"
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err1 := sql.Open("postgres", connStr)

	if err1 != nil {
		//log.Fatal("This is the error: ", err1)
		fmt.Printf("Cannot connect to %s database", connStr)
		//return nil
		panic(err1)
	}

	var err3 = db.Ping()
	if err3 != nil {
		//log.Fatal(err3)
		fmt.Printf("Cannot ping %s database", connStr)
		panic(err3)
		//return nil
	}

	create_table_query := `CREATE TABLE IF NOT EXISTS ways
	(
		id SERIAL,
		charid varchar(50) NOT NULL,
		locationname varchar(50) NOT NULL,
		localetime varchar(50) NOT NULL,
		PRIMARY KEY (id)
	)`

	_, err2 := db.Exec(create_table_query)
	if err2 != nil {
		panic(err2)
	}

	return db
}

func GetWays(w http.ResponseWriter, r *http.Request) {

	CharID := r.FormValue("char_id")

	var response = JSONresponse{}

	if CharID == "" {
		response = JSONresponse{Type: "error", Message: "You are missing one of parameters."}
	} else {
		db := setupDB()

		rows, err := db.Query(`SELECT * FROM ways where charid = $1`, CharID)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		way := []Ways{}

		for rows.Next() {
			p := Ways{}
			err := rows.Scan(&p.id, &p.char_id, &p.location, &p.the_time)
			if err != nil {
				fmt.Println(err)
				continue
			}
			way = append(way, p)
		}

		response = JSONresponse{Type: "success", Data: way}
		defer db.Close()
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteWays(w http.ResponseWriter, r *http.Request) {

	CharID := r.FormValue("char_id")

	var response = JSONresponse{}

	if CharID == "" {
		response = JSONresponse{Type: "error", Message: "You are missing one of parameters."}
	} else {
		db := setupDB()
		_, err := db.Exec(`DELETE FROM ways where charid = $1`, CharID)
		if err != nil {
			panic(err)
		}
		response = JSONresponse{Type: "success", Message: "All ways belong char_id have been deleted successfully!"}
		defer db.Close()
	}
	json.NewEncoder(w).Encode(response)
}

func CreateWays(w http.ResponseWriter, r *http.Request) {

	CharID := r.FormValue("char_id")
	LocationName := r.FormValue("location")
	LocaleTime := r.FormValue("the_time")

	fmt.Println(CharID + " " + LocaleTime + " " + LocaleTime)

	var response = JSONresponse{}

	if CharID == "" || LocationName == "" || LocaleTime == "" {
		response = JSONresponse{Type: "error", Message: "You are missing one of parameters."}
	} else {
		db := setupDB()

		var lastInsertID int
		err := db.QueryRow(`INSERT INTO ways (charid, locationname, localetime) values ($1, $2, $3) returning id;`,
			CharID, LocationName, LocaleTime).Scan(&lastInsertID)
		if err != nil {
			panic(err)
		}
		response = JSONresponse{Type: "success", Message: "The way has been inserted successfully!"}
		defer db.Close()
	}

	json.NewEncoder(w).Encode(response)
}

func main() {

	router := mux.NewRouter()

	// Get all ways by char_id
	router.HandleFunc("/way/", GetWays).Methods("GET")

	// Create a way by char_id
	router.HandleFunc("/way/", CreateWays).Methods("POST")

	// Delete all ways by char_id
	router.HandleFunc("/way/", DeleteWays).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 10000")
	log.Fatal(http.ListenAndServe(":10000", router))

}
