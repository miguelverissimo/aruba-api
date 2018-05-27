package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

type BoatHasSkipair struct {
	ID                uint `gorm:"primary_key"`
	IDBoat            uint `gorm:"foreignkey:id_boat"`
	IDStudentSkipair1 uint `gorm:"foreignkey:id_student_skipair1"`
	IDStudentSkipair2 uint `gorm:"foreignkey:id_student_skipair2"`
	IDStudentSkipair3 uint `gorm:"foreignkey:id_student_skipair3"`
	IDStudentSkipair4 uint `gorm:"foreignkey:id_student_skipair4"`
}

// needed because of poor table naming
func (BoatHasSkipair) TableName() string {
	return "boat_has_skipair"
}

func (i *Server) GetAllSkisOnBoat(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	skiPairs := BoatHasSkipair{}
	if i.DB.Where("id_boat = ?", id).Take(&skiPairs).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&skiPairs)
}

func (i *Server) PostSkisOnBoat(w rest.ResponseWriter, r *rest.Request) {
	skiPairs := BoatHasSkipair{}

	id, err := strconv.ParseUint(r.PathParam("id"), 10, 64)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&skiPairs); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	skiPairs.IDBoat = uint(id)
	if err := i.DB.Save(&skiPairs).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&skiPairs)
}

func (i *Server) PutSkisOnBoat(w rest.ResponseWriter, r *rest.Request) {
	skiPairs := BoatHasSkipair{}

	id := r.PathParam("id")
	if i.DB.Where("id_boat = ?", id).Take(&skiPairs).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := BoatHasSkipair{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	skiPairs.IDStudentSkipair1 = updated.IDStudentSkipair1
	skiPairs.IDStudentSkipair2 = updated.IDStudentSkipair2
	skiPairs.IDStudentSkipair3 = updated.IDStudentSkipair3
	skiPairs.IDStudentSkipair4 = updated.IDStudentSkipair4

	if err := i.DB.Save(&skiPairs).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&skiPairs)

}

func (i *Server) DeleteAllSkisOnBoat(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	skiPairs := BoatHasSkipair{}
	if i.DB.Where("id_boat = ?", id).Take(&skiPairs).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&skiPairs).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
