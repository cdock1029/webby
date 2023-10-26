package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	webbyDb "github.com/cdock1029/webby/db"

	"database/sql"

	"github.com/lib/pq"

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

		t.ExecuteTemplate(w, "index.html.tmpl", map[string]any{"Properties": properties})
	})
	r.Get("/properties", func(w http.ResponseWriter, r *http.Request) {

		properties, err := repo.ListProperties(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t.ExecuteTemplate(w, "_fragment.properties.html.tmpl", map[string]any{"Properties": properties})
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
			data := make(map[string]any)
			if err, ok := err.(*pq.Error); ok {
				log.Default().Println("Database error: ", err.Message)
				log.Default().Println("Database error: ", err.Detail)
				log.Default().Println("Database error: ", err.Hint)
				data["Error"] = err.Detail
			}
			t.ExecuteTemplate(w, "_fragment.property.form.html.tmpl", data)
			return
		}
		w.Header().Set("HX-Trigger", "newProperty")
		w.WriteHeader(204)
	})
	r.Get("/property/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
		edit := r.URL.Query().Get("edit")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		property, err := repo.GetProperty(r.Context(), int32(id))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if edit != "" {
			t.ExecuteTemplate(w, "_fragment.property.edit.html.tmpl", property)
		} else {
			t.ExecuteTemplate(w, "_fragment.property.html.tmpl", property)
		}
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
		w.Header().Set("HX-Trigger", fmt.Sprintf("propertyDeleted%v", id))
		//w.WriteHeader(200)
		w.WriteHeader(204)
	})
	r.Put("/property/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			return
		}
		form := r.PostForm
		name := strings.TrimSpace(form.Get("name"))
		if name == "" {
			w.WriteHeader(204)
			return
		}
		property, err := repo.UpdateProperty(r.Context(), webbyDb.UpdatePropertyParams{ID: int32(id), Name: name})
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("HX-Trigger", fmt.Sprintf("propertyUpdated%v", property.ID))
		t.ExecuteTemplate(w, "_fragment.property.html.tmpl", property)
	})
	r.Get("/null", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}
