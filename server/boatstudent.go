package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

type BoatHasStudent struct {
	ID        uint `gorm:"primary_key"`
	IDBoat    uint `gorm:"column:id_boat"`
	IDStudent uint `gorm:"column:id_student"`
}

// needed because of poor table naming
func (BoatHasStudent) TableName() string {
	return "boat_has_student"
}

func (i *Server) GetBoatForStudent(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	boatForStudent := BoatHasStudent{}
	if i.DB.Where("id_student = ?", id).Take(&boatForStudent).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&boatForStudent)
}

func (i *Server) PostStudentOnBoat(w rest.ResponseWriter, r *rest.Request) {
	boatForStudent := BoatHasStudent{}

	id, err := strconv.ParseUint(r.PathParam("id"), 10, 64)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&boatForStudent); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boatForStudent.IDStudent = uint(id)
	if err := i.DB.Save(&boatForStudent).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&boatForStudent)
}

func (i *Server) PutStudentOnBoat(w rest.ResponseWriter, r *rest.Request) {
	boatForStudent := BoatHasStudent{}

	id := r.PathParam("id")
	if i.DB.Where("id_student = ?", id).Take(&boatForStudent).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := BoatHasStudent{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boatForStudent.IDBoat = updated.IDBoat

	if err := i.DB.Save(&boatForStudent).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&boatForStudent)
}

func (i *Server) DeleteBoatForStudent(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	boatForStudent := BoatHasStudent{}
	if i.DB.Where("id_student = ?", id).Take(&boatForStudent).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&boatForStudent).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
