package dto

type UserReservationsDto struct {
	UserId       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	UserLastName string `json:"user_last_name"`
	UserDni      string `json:"user_dni"`
	UserEmail    string `json:"user_email"`
	UserPassword string `json:"user_password"`

	Reservations ReservationsDto `json:"reservations,omitempty"`
}
