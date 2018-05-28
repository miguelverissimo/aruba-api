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

type StudentFull struct {
	Student
	BoatID     uint
	SkisBoatID uint
}

// needed because of poor table naming
func (Student) TableName() string {
	return "student"
}

func (i *Server) GetAllStudents(w rest.ResponseWriter, r *rest.Request) {
	students := []StudentFull{}
	i.DB.Find(&students)

	i.DB.Debug().Raw(`
			SELECT
				s.id, s.first_name, s.last_name, s.has_skipair,
				bs.id_boat as boat_id, skis.id_boat as skis_boat_id
			FROM student s
			LEFT JOIN boat_has_student bs ON (bs.id_student = s.id)
			LEFT JOIN boat_has_skipair skis 
					ON (skis.id_student_skipair1 = s.id 
					OR skis.id_student_skipair2 = s.id 
					OR skis.id_student_skipair3 = s.id 
					OR skis.id_student_skipair4 = s.id)`).Scan(&students)

	w.WriteJson(&students)
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
