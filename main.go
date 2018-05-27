package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/miguelverissimo/aruba-api/server"
	"log"
	"net/http"
)

func main() {

	s := server.Server{}
	s.InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/boats", s.GetAllBoats),
		rest.Post("/boats", s.PostBoat),
		rest.Get("/boats/:id", s.GetBoat),
		rest.Put("/boats/:id", s.PutBoat),
		rest.Delete("/boats/:id", s.DeleteBoat),

		rest.Get("/students", s.GetAllStudents),
		rest.Post("/students", s.PostStudent),
		rest.Get("/students/:id", s.GetStudent),
		rest.Put("/students/:id", s.PutStudent),
		rest.Delete("/students/:id", s.DeleteStudent),

		rest.Get("/books", s.GetAllBooks),
		rest.Post("/books", s.PostBook),
		rest.Get("/books/:id", s.GetBook),
		rest.Put("/books/:id", s.PutBook),
		rest.Delete("/books/:id", s.DeleteBook),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
