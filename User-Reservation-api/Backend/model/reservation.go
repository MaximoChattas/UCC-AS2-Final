package model

type Reservation struct {
	Id        int     `gorm:"primaryKey"`
	StartDate string  `gorm:"type:varchar(16); not null"` //Expected time as "DD-MM-YYYY hh:mm"
	EndDate   string  `gorm:"type:varchar(16); not null"` //Expected time as "DD-MM-YYYY hh:mm"
	UserId    int     `gorm:"foreignkey:UserId"`
	HotelId   string  `gorm:"type:varchar(24); not null"`
	Amount    float64 `gorm:"type:decimal(10,2); not null"`
}

type Reservations []Reservation
