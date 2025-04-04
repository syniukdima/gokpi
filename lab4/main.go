package main

import (
	"Lab4Web/handlers"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/task1", handlers.Task1Handler)
	http.HandleFunc("/task2", handlers.Task2Handler)
	http.HandleFunc("/task3", handlers.Task3Handler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/index.html"))
	tmpl.ExecuteTemplate(w, "layout", nil)
}
