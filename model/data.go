package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	ActivityDate time.Time `json:"activity_date"`
	LinkMom      string    `json:"link_mom"`
	Location     string    `json:"location"`
	Actor        string    `json:"actor"`
}

type Datas struct {
	Data []Data `json:"data"`
}

func (data *Data) BeforeCreate(tx *gorm.DB) (err error) {
	data.ID = uuid.New()
	return
}

type DataWithoutFunc struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	ActivityDate time.Time `json:"activity_date"`
	LinkMom      string    `json:"link_mom"`
	Location     string    `json:"location"`
	Actor        string    `json:"actor"`
}

func (data Data) MarshalJSON() ([]byte, error) {
	dataWithoutFunc := DataWithoutFunc{
		ID:           data.ID,
		Title:        data.Title,
		Description:  data.Description,
		Image:        data.Image,
		ActivityDate: data.ActivityDate,
		LinkMom:      data.LinkMom,
		Location:     data.Location,
		Actor:        data.Actor,
	}

	return json.Marshal(&struct {
		DataWithoutFunc
	}{
		DataWithoutFunc: dataWithoutFunc,
	})
}
