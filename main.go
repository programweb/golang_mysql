package main

import (
	"database/sql"

	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type CauseHealthData struct {
	Id    int
	Short string
	Name  string
}

// cooler connection stuff here: https://cloud.google.com/sql/docs/mysql/samples/cloud-sql-mysql-databasesql-connect-tcp

func dbConn() (db *sql.DB) {
	//dbDriver := "mysql"
	//dbUser := "root"
	//dbPass := "ABC"
	//dbName := "tcp(127.0.0.1:33061)/goblog"
	// db, err := sql.Open("mysql", "root:ABC@tcp(127.0.0.1:33061)/healthdata")
	// db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	db, err := sql.Open("mysql", "root:ABC@tcp(127.0.0.1:33061)/healthdata")
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))
var fserve = http.StripPrefix("/", http.FileServer(http.Dir("static")))

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fserve.ServeHTTP(w, r)
		return
	}
	db := dbConn()
	selDB, err := db.Query("SELECT cause_id, cause_short, cause_name FROM healthdata.cause ORDER BY cause_id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := CauseHealthData{}
	res := []CauseHealthData{}
	for selDB.Next() {
		var cause_id int
		var cause_short, cause_name string
		err = selDB.Scan(&cause_id, &cause_short, &cause_name)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = cause_id
		emp.Short = cause_short
		emp.Name = cause_name
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT cause_id, cause_short, cause_name FROM cause WHERE cause_id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := CauseHealthData{}
	for selDB.Next() {
		var cause_id int
		var cause_short, cause_name string
		err = selDB.Scan(&cause_id, &cause_short, &cause_name)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = cause_id
		emp.Short = cause_short
		emp.Name = cause_name
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT cause_id, cause_short, cause_name FROM cause WHERE cause_id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := CauseHealthData{}
	for selDB.Next() {
		var cause_id int
		var cause_short, cause_name string
		err = selDB.Scan(&cause_id, &cause_short, &cause_name)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = cause_id
		emp.Short = cause_short
		emp.Name = cause_name
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		key := r.FormValue("key")
		short := r.FormValue("short")
		name := r.FormValue("name")
		insForm, err := db.Prepare("INSERT INTO cause(acause, cause_short, cause_name) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(key, short, name)
		log.Println("INSERT:  Key (acause): " + key + " | Short Name: " + short + " | Name: " + name)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		short := r.FormValue("short")
		name := r.FormValue("name")
		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE cause SET cause_short=?, cause_name=? WHERE cause_id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(short, name, id)
		log.Println("UPDATE: Short: " + name + " | Name: " + name)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM cause WHERE cause_id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Golang Server started on: ")
	log.Println("Server web application at: http://localhost:8080")

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	http.ListenAndServe(":8080", nil) // ListenAndServe starts an HTTP server
}
