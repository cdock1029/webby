package main

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"

	webby "example.com/webby/sql"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates
var templates embed.FS

//go:embed static
var static embed.FS

var t = template.Must(template.ParseFS(templates, "templates/*"))

func run() error {
	ctx := context.Background()

	db, err := sql.Open("postgres", "host=/run/postgresql dbname=webby")
	if err != nil {
		return err
	}
	queries := webby.New(db)

	properties, err := queries.ListProperties(ctx)
	if err != nil {
		return err
	}
	log.Print(properties)

	// inserted, err := queries.CreateProperty(ctx, gofakeit.Company())
	// if err != nil {
	// 	return err
	// }
	// log.Print(inserted)

	// fetched, err := queries.GetProperty(ctx, inserted.ID)
	// if err != nil {
	// 	return err
	// }
	// log.Println(reflect.DeepEqual(inserted, fetched))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	definemiscellaneousRoutes(r)
	defineAppRoutes(r)

	log.Println("listening on :4000")
	http.ListenAndServe(":4000", r)
}

func defineAppRoutes(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		t.ExecuteTemplate(w, "index.html.tmpl", [0]string{})
	})

}

func definemiscellaneousRoutes(r *chi.Mux) {
	r.Use(middleware.Logger)
	fs := http.FileServer(http.FS(static))
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})
	r.Handle("/static/*", fs)
}
