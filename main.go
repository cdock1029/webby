package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates
var templates embed.FS

//go:embed static
var static embed.FS

var t = template.Must(template.ParseFS(templates, "templates/*"))

type Todo struct {
	Title string
	note  string
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fs := http.FileServer(http.FS(static))

	r.Get("/favicon.ico", serveFavIcon)
	r.Handle("/static/*", fs)

	r.Get("/", getHome)

	log.Println("listening on :4000")
	http.ListenAndServe(":4000", r)
}

func serveFavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func getHome(w http.ResponseWriter, r *http.Request) {
	todos := []Todo{
		{Title: "hello world", note: "first todo"},
		{Title: "another todo"},
		{Title: "making a list", note: "yes"},
	}

	//t.ExecuteTemplate(w, "index.html.tmpl", make(map[string]string))
	t.ExecuteTemplate(w, "index.html.tmpl", todos)
}
