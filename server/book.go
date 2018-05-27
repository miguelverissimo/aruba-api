package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type Book struct {
	ID uint `gorm:"primary_key"`

	Name        string `gorm:"size:255"`
	URLOnAmazon string `gorm:"size:255"`
}

// needed because of poor table naming
func (Book) TableName() string {
	return "book"
}

func (i *Server) GetAllBooks(w rest.ResponseWriter, r *rest.Request) {
	books := []Book{}
	i.DB.Find(&books)
	w.WriteJson(&books)
}

func (i *Server) GetBook(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	book := Book{}
	if i.DB.First(&book, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&book)
}

func (i *Server) PostBook(w rest.ResponseWriter, r *rest.Request) {
	book := Book{}
	if err := r.DecodeJsonPayload(&book); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&book).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&book)
}

func (i *Server) PutBook(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	book := Book{}
	if i.DB.First(&book, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := Book{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.Name = updated.Name
	book.URLOnAmazon = updated.URLOnAmazon

	if err := i.DB.Save(&book).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&book)
}

func (i *Server) DeleteBook(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	book := Book{}
	if i.DB.First(&book, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&book).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
