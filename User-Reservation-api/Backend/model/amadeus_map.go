package model

type AmadeusMap struct {
	HotelId   string `gorm:"type:varchar(24); primaryKey"`
	AmadeusId string `gorm:"type:varchar(8); not null"`
}
