package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"

	webbyDb "github.com/cdock1029/webby/db"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates/*
var templates embed.FS

//go:embed static
var static embed.FS

var t = template.Must(template.ParseFS(templates, "templates/*"))

var repo *webbyDb.Queries

func setupDb() error {
	sqlDb, err := sql.Open("postgres", "host=/run/postgresql dbname=webby")
	if err != nil {
		return err
	}
	repo = webbyDb.New(sqlDb)
	return nil
}

func main() {
	if err := setupDb(); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/static/*", http.FileServer(http.FS(static)))
	defineAppRoutes(r)

	log.Println("listening on :4000")
	http.ListenAndServe(":4000", r)
}

func defineAppRoutes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		properties, err := repo.ListProperties(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t.ExecuteTemplate(w, "index.html.tmpl", properties)
	})
	r.Get("/properties", func(w http.ResponseWriter, r *http.Request) {

		properties, err := repo.ListProperties(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t.ExecuteTemplate(w, "_fragment.properties.html.tmpl", properties)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			return
		}
		property := r.PostForm
		name := property.Get("name")
		if name == "" {
			t.ExecuteTemplate(w, "error.html.tmpl", "Name can't be blank")
			return
		}
		if _, err := repo.CreateProperty(r.Context(), name); err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("HX-Trigger", "newProperty")
		w.WriteHeader(204)
	})
	r.Delete("/property/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		err = repo.DeleteProperty(r.Context(), int32(id))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})
}
