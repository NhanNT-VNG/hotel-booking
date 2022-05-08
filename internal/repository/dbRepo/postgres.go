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
			room_id = $1 and
			$2 < end_date and $3 > start_date;
	`
	row := m.DB.QueryRowContext(ctx, query, roomId, statDate, endDate)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

func (m *postgresDBRepo) SearchAvailabilityAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select r.id, r.room_name
		from rooms r
		where r.id not in 
		(select room_id from room_restrictions where $1 < end_date and $2 > start_date);
	`
	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)

	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (m *postgresDBRepo) GetRoomById(roomId int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, roomId)

	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}
