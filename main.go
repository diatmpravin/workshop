package main

import (
	"fmt"
	"net/http"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"text/template"
)

type User struct{
	Fname string
	Email string
	Phone string
	Country string
	Coursetype string
}

func setupDB() *sql.DB{
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/registration?charset=utf8")
	PanicIf(err)
	return db
}

func PanicIf(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

var (
	db *sql.DB
	createTable = `CREATE TABLE IF NOT EXISTS user (
		fname VARCHAR(64) NULL DEFAULT NULL,
		email VARCHAR(64) NULL DEFAULT NULL,
		phone VARCHAR(64) NULL DEFAULT NULL,
		country VARCHAR(64) NULL DEFAULT NULL
    );`
)

func newHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("New handler called")
	t := template.New("new.html")
	t.ParseFiles("templates/new.html")
	t.Execute(w, t)
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("View handler called")
	rows, err := db.Query("Select * from user")
	PanicIf(err)
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Fname, &user.Email, &user.Phone, &user.Country)
		PanicIf(err)
		users = append(users, user)
	}
	t := template.New("new.html")
	t, _ = template.ParseFiles("templates/index.html")
	t.Execute(w, users)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Save handler called")
	stmt, err := db.Prepare("INSERT user SET fname=?,email=?,phone=?,country=?")
	PanicIf(err)
	fname := r.FormValue("fname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	country := r.FormValue("country")
	res, err := stmt.Exec(fname, email, phone, country)
	PanicIf(err)
	fmt.Println(res)
	http.Redirect(w, r, "/view/" , http.StatusFound)
}


func main() {
	db = setupDB()
	defer db.Close()

	ctble, err := db.Query(createTable)
	PanicIf(err)
	fmt.Println("Table create successull", ctble)


	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/save/", saveHandler)

	http.Handle("/public/css/", http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/public/images/", http.StripPrefix("/public/images/", http.FileServer(http.Dir("public/images"))))

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	fmt.Println("Listening Server.....")
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Template Execution %s", err)
	}
}

