package dto

type AmadeusMapDto struct {
	HotelId   string `json:"hotel_id" validate:"required"`
	AmadeusId string `json:"amadeus_id" validate:"required"`
}
