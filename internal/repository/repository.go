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
}
