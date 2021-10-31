package main

import (
	"fmt"
	"net/http"
	"log"
	"text/template"
)

func newHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("New handler called")
	t := template.New("new.html")
	t.ParseFiles("templates/new.html")
	t.Execute(w, t)
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("View handler called")
	t := template.New("new.html")
	t, _ = template.ParseFiles("templates/index.html")
	t.Execute(w, t)
}


func main() {

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/new/", newHandler)

	http.Handle("/public/css/", http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/public/images/", http.StripPrefix("/public/images/", http.FileServer(http.Dir("public/images"))))

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	fmt.Println("Listening Server.....")
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Template Execution %s", err)
	}
}

