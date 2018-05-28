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
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			// return origin == "http://aruba-client.127.0.0.1.xip.io/"
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlExposeHeaders: []string{
			"Access-Control-Allow-Origin", "*"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	router, err := rest.MakeRouter(
		rest.Get("/boats", s.GetAllBoats),
		rest.Post("/boats", s.PostBoat),
		rest.Get("/boats/:id", s.GetBoat),
		rest.Put("/boats/:id", s.PutBoat),
		rest.Delete("/boats/:id", s.DeleteBoat),

		rest.Get("/boats/:id/skis", s.GetAllSkisOnBoat),
		rest.Post("/boats/:id/skis", s.PostSkisOnBoat),
		rest.Put("/boats/:id/skis", s.PutSkisOnBoat),
		rest.Delete("/boats/:id/skis", s.DeleteAllSkisOnBoat),

		rest.Get("/students", s.GetAllStudents),
		rest.Post("/students", s.PostStudent),
		rest.Get("/students/:id", s.GetStudent),
		rest.Put("/students/:id", s.PutStudent),
		rest.Delete("/students/:id", s.DeleteStudent),

		rest.Get("/students/:id/boat", s.GetBoatForStudent),
		rest.Post("/students/:id/boat", s.PostStudentOnBoat),
		rest.Put("/students/:id/boat", s.PutStudentOnBoat),
		rest.Delete("/students/:id/boat", s.DeleteBoatForStudent),

		rest.Get("/books", s.GetAllBooks),
		rest.Post("/books", s.PostBook),
		rest.Get("/books/:id", s.GetBook),
		rest.Put("/books/:id", s.PutBook),
		rest.Delete("/books/:id", s.DeleteBook),

		rest.Get("/books/:id/boat", s.GetBoatForBook),
		rest.Post("/books/:id/boat", s.PostBookOnBoat),
		rest.Put("/books/:id/boat", s.PutBookOnBoat),
		rest.Delete("/books/:id/boat", s.DeleteBoatForBook),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
