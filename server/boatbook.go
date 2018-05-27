package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

type BoatHasBook struct {
	ID     uint `gorm:"primary_key"`
	IDBoat uint `gorm:"column:id_boat"`
	IDBook uint `gorm:"column:id_book"`
}

// needed because of poor table naming
func (BoatHasBook) TableName() string {
	return "boat_has_book"
}

func (i *Server) GetBoatForBook(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	boatForBook := BoatHasBook{}
	if i.DB.Where("id_book = ?", id).Take(&boatForBook).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&boatForBook)
}

func (i *Server) PostBookOnBoat(w rest.ResponseWriter, r *rest.Request) {
	boatForBook := BoatHasBook{}

	id, err := strconv.ParseUint(r.PathParam("id"), 10, 64)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&boatForBook); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boatForBook.IDBook = uint(id)
	if err := i.DB.Save(&boatForBook).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&boatForBook)
}

func (i *Server) PutBookOnBoat(w rest.ResponseWriter, r *rest.Request) {
	boatForBook := BoatHasBook{}

	id := r.PathParam("id")
	if i.DB.Where("id_book = ?", id).Take(&boatForBook).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := BoatHasBook{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boatForBook.IDBoat = updated.IDBoat

	if err := i.DB.Save(&boatForBook).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&boatForBook)
}

func (i *Server) DeleteBoatForBook(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	boatForBook := BoatHasBook{}
	if i.DB.Where("id_book = ?", id).Take(&boatForBook).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&boatForBook).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
