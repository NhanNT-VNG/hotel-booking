package repository

import (
	"time"

	"github.com/NhanNT-VNG/hotel-booking/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestrictions(rr models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(statDate, endDate time.Time, roomId int) (bool, error)
	SearchAvailabilityAllRooms(startDate, endDate time.Time) ([]models.Room, error)
	GetRoomById(roomId int) (models.Room, error)

	GetUserById(userId int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, password string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)
	UpdateReservation(reservation models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedReservation(id, processed int) error
}
