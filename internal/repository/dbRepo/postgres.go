package dbrepo

import (
	"context"
	"time"

	"github.com/NhanNT-VNG/hotel-booking/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(reservation models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into reservations(
		first_name, last_name, email, phone, start_date, 
		end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var reservationId int

	err := m.DB.QueryRowContext(ctx, query,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		reservation.StartDate,
		reservation.EndDate,
		reservation.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&reservationId)

	if err != nil {
		return 0, err
	}
	return reservationId, nil
}

func (m *postgresDBRepo) InsertRoomRestrictions(rr models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into room_restrictions(
		start_date, end_date, room_id, reservation_id, 
		restriction_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(
		ctx, query,
		rr.StartDate,
		rr.EndDate,
		rr.RoomId,
		rr.ReservationId,
		rr.RestrictionId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomId(statDate, endDate time.Time, roomId int) (bool, error) {
	ctx, cancer := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancer()
	var numRows int
	query := `
		select count(id)
		from room_restrictions
		where 
			room_id = $1
			$2 < end_date and $3 > start_date;
	`
	row := m.DB.QueryRowContext(ctx, query, roomId, statDate, endDate)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return false, nil
	}
	return true, nil
}
