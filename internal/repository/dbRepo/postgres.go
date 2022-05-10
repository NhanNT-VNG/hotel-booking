package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/NhanNT-VNG/hotel-booking/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func (m *postgresDBRepo) GetUserById(userId int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select 
		id, first_name, last_name, email, password, access_level, created_at, updated_at
	from users
	where id = $1`
	row := m.DB.QueryRowContext(ctx, query, userId)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}
	return user, nil
}

func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		update users
		set 
			first_name = $1,
			last_name = $2,
			email = $3,
			access_level = $4,
			updated_at = $5
		where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect email or password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservationList []models.Reservation

	query := `
		select 
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
			r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name,
			r.processed
		from reservations r
		left join rooms rm on rm.id = r.room_id
		order by r.start_date
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservationList, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomId,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
			&reservation.Processed,
		)

		if err != nil {
			return reservationList, err
		}
		reservationList = append(reservationList, reservation)
	}

	if err = rows.Err(); err != nil {
		return reservationList, err
	}
	return reservationList, nil
}

func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservationList []models.Reservation

	query := `
		select 
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
			r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name
		from reservations r
		left join rooms rm on rm.id = r.room_id
		where r.processed = 0
		order by r.start_date
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservationList, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomId,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)

		if err != nil {
			return reservationList, err
		}
		reservationList = append(reservationList, reservation)
	}

	if err = rows.Err(); err != nil {
		return reservationList, err
	}
	return reservationList, nil
}

func (m *postgresDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	query := `
		select 
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
			r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name,
			r.processed
		from reservations r
		left join rooms rm on rm.id = r.room_id
		where r.id = $1
	`
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomId,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
		&reservation.Processed,
	)

	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

func (m *postgresDBRepo) UpdateReservation(reservation models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		update reservations
		set 
			first_name = $1,
			last_name = $2,
			email = $3,
			phone = $4,
			updated_at = $5
		where id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		time.Now(),
		reservation.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from reservations where id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) UpdateProcessedReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update reservations set processed = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `select id, room_name, created_at, updated_at from rooms order by room_name`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {
		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
		)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Close(); err != nil {
		return rooms, err
	}

	return rooms, nil
}
