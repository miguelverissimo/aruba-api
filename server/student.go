package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type Student struct {
	ID         uint   `gorm:"primary_key"`
	FirstName  string `gorm:"size:255"`
	LastName   string `gorm:"size:255"`
	HasSkipair bool
}

// needed because of poor table naming
func (Student) TableName() string {
	return "student"
}

func (i *Server) GetAllStudents(w rest.ResponseWriter, r *rest.Request) {
	student := []Student{}
	i.DB.Find(&student)
	w.WriteJson(&student)
}

func (i *Server) GetStudent(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	student := Student{}
	if i.DB.First(&student, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&student)
}

func (i *Server) PostStudent(w rest.ResponseWriter, r *rest.Request) {
	student := Student{}
	if err := r.DecodeJsonPayload(&student); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&student).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&student)
}

func (i *Server) PutStudent(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	student := Student{}
	if i.DB.First(&student, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := Student{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	student.FirstName = updated.FirstName
	student.LastName = updated.LastName
	student.HasSkipair = updated.HasSkipair

	if err := i.DB.Save(&student).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&student)
}

func (i *Server) DeleteStudent(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	student := Student{}
	if i.DB.First(&student, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&student).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
