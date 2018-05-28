package server

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"strings"
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

var COLOR_ENUM = map[string]int{
	"BLUE":      1,
	"NAVY_BLUE": 2,
	"GREEN":     3,
	"RED":       4,
	"PURPLE":    5,
	"WHITE":     6,
	"BLACK":     7,
	"YELLOW":    8,
}

type Boat struct {
	ID                uint   `gorm:"primary_key"`
	Name              string `gorm:"size:255"`
	Price             float32
	Color             BoatColor `sql:"enum('BLUE','NAVY_BLUE','GREEN','RED','PURPLE','WHITE','BLACK','YELLOW')"`
	LastStudentChange time.Time `sql:"DEFAULT: NULL"`
}

type BoatFull struct {
	Boat
	Students string
	Books    string
	Skis     string
}

type BoatFullResponse struct {
	Boat
	Students []uint
	Books    []uint
	Skis     []uint
}

// TableName is needed because of poor table naming (should be "boats")
func (Boat) TableName() string {
	return "boat"
}

func IDsStringToUintSlice(str string) []uint {
	var ids []uint
	if len(str) == 0 {
		return ids
	}
	s := strings.Split(str, ",")
	for _, i := range s {
		if val, err := strconv.Atoi(i); err == nil {
			ids = append(ids, uint(val))
		}
	}
	return ids
}

func (b *BoatFull) ConvertToBoatFullResponse() BoatFullResponse {
	if b == nil {
		return BoatFullResponse{}
	}

	return BoatFullResponse{
		Boat: Boat{
			ID:    b.ID,
			Name:  b.Name,
			Price: b.Price,
			Color: b.Color,
		},
		Students: IDsStringToUintSlice(b.Students),
		Books:    IDsStringToUintSlice(b.Books),
		Skis:     IDsStringToUintSlice(b.Skis),
	}
}

func (u *BoatColor) Scan(value interface{}) error {
	*u = BLUE
	if v, ok := value.([]uint8); ok {
		val := string(v)
		if COLOR_ENUM[val] >= 1 && COLOR_ENUM[val] <= 8 {
			fmt.Printf("val: %v, COLOR_ENUM[val]: %v, BoatColor(COLOR_ENUM[val]): %v", val, COLOR_ENUM[val], BoatColor(COLOR_ENUM[val]))
			*u = BoatColor(COLOR_ENUM[val])
		} else {
			fmt.Printf("Still wrong inside:\n%T: %v\n\n%T: %v\n", value, value, val, val)
		}
	} else {
		fmt.Printf("Still wrong:\n%T: %+v\n", value, value)
	}
	return nil
}

func (i *Server) GetAllBoats(w rest.ResponseWriter, r *rest.Request) {
	boats := []BoatFull{}
	boatsResponse := []BoatFullResponse{}

	i.DB.Debug().Raw(`
			SELECT
				b.id, b.name, b.price, b.color, b.last_student_change,
				GROUP_CONCAT(DISTINCT id_student ORDER BY id_student SEPARATOR ',') students,
				GROUP_CONCAT(DISTINCT id_book ORDER BY id_book SEPARATOR ',') books,
				GROUP_CONCAT(DISTINCT id_student ORDER BY id_student SEPARATOR ',') skis
			FROM boat b
			LEFT JOIN      boat_has_student bs ON (bs.id_boat = b.id)
			LEFT JOIN      boat_has_book bb ON (bb.id_boat = b.id)
			LEFT JOIN      boat_has_skipair bsp ON (bsp.id_boat = b.id)
			GROUP BY  b.id`).Scan(&boats)
	for _, boat := range boats {
		boatsResponse = append(boatsResponse, boat.ConvertToBoatFullResponse())
	}
	w.WriteJson(&boatsResponse)
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
