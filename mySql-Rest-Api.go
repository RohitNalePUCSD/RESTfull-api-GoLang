package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	DBUSER string = os.Getenv("DATABASE_USERNAME")
	DBPASS string = os.Getenv("DATABASE_PASSWORD")
	DBNAME string = os.Getenv("DATABASE_NAME")
)

type articals struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Descp   string `json:"desc"`
	Content string `json:"content"`
}

func InitRouter() (router *mux.Router) {

	router = mux.NewRouter()

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/artical", getPosts).Methods("GET")
	router.HandleFunc("/artical", postPost).Methods("POST")
	router.HandleFunc("/artical/{id}", getPost).Methods("GET")
	router.HandleFunc("/artical/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/artical/{id}", deletePost).Methods("DELETE")

	return
}

func dbConnection() (conn *sql.DB) {

	conn, err := sql.Open("mysql", "rohit:rohit@/articalDB")
	if err != nil {
		panic(err)
	}
	return
}

func postPost(w http.ResponseWriter, r *http.Request) {

	db := dbConnection()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	//fmt.Fprintf(w, "%+v", string(body))

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	Id := keyVal["Id"]
	Title := keyVal["Title"]
	Descp := keyVal["Descp"]
	Content := keyVal["Content"]

	fmt.Println("keyvalue----------->", Id, Title, Descp, Content)
	result, err := db.Query("INSERT INTO artical VALUES (?, ?, ?, ?)", Id, Title, Descp, Content)
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	fmt.Fprintf(w, "artical with ID = %s was inserted", keyVal["id"])

}

func deletePost(w http.ResponseWriter, r *http.Request) {

	db := dbConnection()
	params := mux.Vars(r)
	result, err := db.Query("DELETE FROM artical WHERE Id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	fmt.Fprintf(w, "artical with ID = %s was deleted", params["id"])

}

func updatePost(w http.ResponseWriter, r *http.Request) {

	db := dbConnection()
	params := mux.Vars(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["Title"]

	result, err := db.Query("UPDATE artical SET Title = ? WHERE Id = ?", newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	fmt.Fprintf(w, "Artical with ID = %s was updated", params["id"])

}

func getPost(w http.ResponseWriter, r *http.Request) {

	db := dbConnection()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT Id, Title, Descp, Content FROM artical WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var artical articals
	for result.Next() {
		err := result.Scan(&artical.Id, &artical.Title, &artical.Descp, &artical.Content)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(artical)
}

func getPosts(w http.ResponseWriter, r *http.Request) {

	var artical []articals
	db := dbConnection()
	query := "select Id, Title, Descp, Content from artical"

	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	for result.Next() {
		var art articals
		err = result.Scan(&art.Id, &art.Title, &art.Descp, &art.Content)
		if err != nil {
			panic(err)
		}
		artical = append(artical, art)
	}
	json.NewEncoder(w).Encode(artical)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to the HomePage Rohit...!")
	fmt.Println("Endpoint Hit: homePage")
}

func serverStart() {

	router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":3000")
}

func main() {

	serverStart()
}
