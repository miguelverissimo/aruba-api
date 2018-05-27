package server

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type BoatColor int

const (
	BLUE      BoatColor = 1
	NAVY_BLUE BoatColor = 2
	GREEN     BoatColor = 3
	RED       BoatColor = 4
	PURPLE    BoatColor = 5
	WHITE     BoatColor = 6
	BLACK     BoatColor = 7
	YELLOW    BoatColor = 8
)

type Boat struct {
	ID                uint   `gorm:"primary_key"`
	Name              string `gorm:"size:255"`
	Price             float32
	Color             BoatColor `sql:"enum('BLUE','NAVY_BLUE','GREEN','RED','PURPLE','WHITE','BLACK','YELLOW')"`
	LastStudentChange time.Time `sql:"DEFAULT: NULL"`
}

// TableName is needed because of poor table naming (should be "boats")
func (Boat) TableName() string {
	return "boat"
}

func (u *BoatColor) Scan(value interface{}) error {
	val, ok := value.(uint)
	if ok {
		*u = BoatColor(val)
	} else {
		*u = BLUE
	}
	return nil
}

// func (u BoatColor) Value() (driver.Value, error) { return string(u), nil }

func (i *Server) GetAllBoats(w rest.ResponseWriter, r *rest.Request) {
	boats := []Boat{}
	i.DB.Find(&boats)
	w.WriteJson(&boats)
}

func (i *Server) GetBoat(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	boat := Boat{}
	if i.DB.First(&boat, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&boat)
}

func (i *Server) PostBoat(w rest.ResponseWriter, r *rest.Request) {
	boat := Boat{}
	if err := r.DecodeJsonPayload(&boat); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&boat).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&boat)
}

func (i *Server) PutBoat(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	boat := Boat{}
	if i.DB.First(&boat, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := Boat{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boat.Name = updated.Name
	boat.Color = updated.Color
	boat.Price = updated.Price
	boat.LastStudentChange = time.Now()

	if err := i.DB.Save(&boat).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&boat)
}

func (i *Server) DeleteBoat(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	boat := Boat{}
	if i.DB.First(&boat, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&boat).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
