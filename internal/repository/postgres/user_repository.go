package repository

import (
	"admin-panel/internal/domain"
	"admin-panel/pkg/lib/utils"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) GetAllUsers(page, pageSize int) (*domain.UsersList, error) {
	offset := (page - 1) * pageSize

	query := `
        SELECT id, first_name, last_name, phone_number, blocked, registration_date, gender, date_of_birth, location, email, profile_photo_url
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		slog.Error("Error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.TODO(), pageSize, offset)
	if err != nil {
		slog.Error("Error executing query: %v", utils.Err(err))
		return nil, err
	}
	defer rows.Close()

	var usersList domain.UsersList

	for rows.Next() {
		var user domain.CommonUserResponse
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Blocked,
			&user.RegistrationDate,
			&user.Gender,
			&user.DateOfBirth,
			&user.Location,
			&user.Email,
			&user.ProfilePhotoURL,
		); err != nil {
			slog.Error("Error scanning user row: %v", utils.Err(err))
			return nil, err
		}

		usersList.Users = append(usersList.Users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over user rows: %v", utils.Err(err))
		return nil, err
	}

	return &usersList, nil
}

func (r *PostgresUserRepository) GetTotalUsersCount() (int, error) {
	var totalUsers int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		slog.Error("error getting total users count", utils.Err(err))
		return 0, err
	}

	return totalUsers, nil
}

func (r *PostgresUserRepository) GetUserByID(id int32) (*domain.GetUserResponse, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, first_name, last_name, phone_number, blocked, registration_date, gender, date_of_birth, location, email, profile_photo_url
		FROM users 
		WHERE id = $1
	`)

	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(context.TODO(), id)

	var user domain.GetUserResponse

	err = row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Blocked,
		&user.RegistrationDate,
		&user.Gender,
		&user.DateOfBirth,
		&user.Location,
		&user.Email,
		&user.ProfilePhotoURL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		slog.Error("error scanning user row: %v", utils.Err(err))
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	if !utils.IsValidPhoneNumber(request.PhoneNumber) {
		return nil, fmt.Errorf("invalid phone number format")
	}

	stmt, err := r.DB.Prepare(`
		INSERT INTO users (first_name, last_name, phone_number,	gender, date_of_birth, location, email, profile_photo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, first_name, last_name, phone_number, blocked,	registration_date, gender, date_of_birth, location,	email, profile_photo_url
	`)
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	var user domain.CreateUserResponse

	err = stmt.QueryRow(
		request.FirstName,
		request.LastName,
		request.PhoneNumber,
		request.Gender,
		request.DateOfBirth,
		request.Location,
		request.Email,
		request.ProfilePhotoURL,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Blocked,
		&user.RegistrationDate,
		&user.Gender,
		&user.DateOfBirth,
		&user.Location,
		&user.Email,
		&user.ProfilePhotoURL,
	)
	if err != nil {
		slog.Error("error executing query: %v", utils.Err(err))
		return nil, err
	}

	return &user, nil
}

func (r PostgresUserRepository) UpdateUser(id int32, request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	updateQuery := `UPDATE users SET
                    first_name = $1,
                    last_name = $2,
                    gender = $3,
                    date_of_birth = $4,
                    location = $5,
                    email = $6,
                    profile_photo_url = $7
                    WHERE id = $8
                    RETURNING id, first_name, last_name, phone_number, blocked, gender, registration_date, date_of_birth, location, email, profile_photo_url`

	stmt, err := r.DB.Prepare(updateQuery)
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	var user domain.UpdateUserResponse
	err = stmt.QueryRow(
		request.FirstName,
		request.LastName,
		request.Gender,
		request.DateOfBirth,
		request.Location,
		request.Email,
		request.ProfilePhotoURL,
		id,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Blocked,
		&user.Gender,
		&user.RegistrationDate,
		&user.DateOfBirth,
		&user.Location,
		&user.Email,
		&user.ProfilePhotoURL,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		slog.Error("error executing query: %v", utils.Err(err))
		return nil, err
	}

	return &user, nil
}

func (r PostgresUserRepository) DeleteUser(id int32) error {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		slog.Error("error checking user existence: %v", utils.Err(err))
		return err
	}

	if !exists {
		return domain.ErrUserNotFound
	}

	stmt, err := r.DB.Prepare(`DELETE FROM users WHERE id = $1`)
	if err != nil {
		slog.Error("error preparing query: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		slog.Error("error executing query: %v", utils.Err(err))
		return err
	}

	return nil
}

func (r *PostgresUserRepository) BlockUser(id int32) error {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		slog.Error("error checking user existence: %v", utils.Err(err))
		return err
	}

	if !exists {
		return domain.ErrUserNotFound
	}

	stmt, err := r.DB.Prepare("UPDATE users SET blocked = true WHERE id = $1")
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		slog.Error("error executing query: %v", utils.Err(err))
		return err
	}

	return nil
}

func (r *PostgresUserRepository) UnblockUser(id int32) error {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		slog.Error("error checking user existence: %v", utils.Err(err))
		return err
	}

	if !exists {
		return domain.ErrUserNotFound
	}

	stmt, err := r.DB.Prepare("UPDATE users SET blocked = false WHERE id = $1")
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		slog.Error("error executing query: %v", utils.Err(err))
		return err
	}

	return nil
}

func (r *PostgresUserRepository) SearchUsers(query string, page, pageSize int) (*domain.UsersList, error) {
	offset := (page - 1) * pageSize

	searchQuery := `
        SELECT id, first_name, last_name, phone_number, blocked,
        registration_date, gender, date_of_birth, location,
        email, profile_photo_url
        FROM users
        WHERE first_name ILIKE $1 OR last_name ILIKE $1 OR phone_number ILIKE $1 OR email ILIKE $1
        ORDER BY id
        LIMIT $2 OFFSET $3
    `

	stmt, err := r.DB.Prepare(searchQuery)
	if err != nil {
		slog.Error("Error preparing search query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.TODO(), "%"+query+"%", pageSize, offset)
	if err != nil {
		slog.Error("Error executing search query: %v", utils.Err(err))
		return nil, err
	}
	defer rows.Close()

	userList := domain.UsersList{Users: make([]domain.CommonUserResponse, 0)}
	for rows.Next() {
		user, err := utils.ScanUserRow(rows)
		if err != nil {
			return nil, err
		}
		userList.Users = append(userList.Users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over user rows: %v", utils.Err(err))
		return nil, err
	}

	return &userList, nil
}
