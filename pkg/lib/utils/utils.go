package utils

import (
	"admin-panel/internal/domain"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func HandleNullString(valid sql.NullString) string {
	if valid.Valid {
		return valid.String
	}
	return ""
}

func NullIfEmptyStr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func NullIfEmptyDate(d domain.Date) sql.NullTime {
	if d == (domain.Date{}) {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true}
}

func IsValidPhoneNumber(phoneNumber string) bool {
	// Check if the phone number consists of 12 digits and starts with "+993"
	const validPrefix = "+993"
	return len(phoneNumber) == 12 && strings.HasPrefix(phoneNumber, validPrefix)
}

// Alternative for http.Error to response with json instead of plain text
func RespondWithErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonError := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  status,
		Message: message,
	}

	json.NewEncoder(w).Encode(jsonError)
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ScanUserRow(rows *sql.Rows) (domain.CommonUserResponse, error) {
	var user domain.CommonUserResponse
	var firstName, lastName, gender, location, email, profilePhotoURL sql.NullString
	var dateOfBirth sql.NullTime

	if err := rows.Scan(
		&user.ID,
		&firstName,
		&lastName,
		&user.PhoneNumber,
		&user.Blocked,
		&user.RegistrationDate,
		&gender,
		&dateOfBirth,
		&location,
		&email, &profilePhotoURL,
	); err != nil {
		slog.Error("Error scanning user row: %v", Err(err))
		return domain.CommonUserResponse{}, err
	}

	user.FirstName = HandleNullString(firstName)
	user.LastName = HandleNullString(lastName)
	user.Gender = HandleNullString(gender)
	user.Location = HandleNullString(location)
	user.Email = HandleNullString(email)
	user.ProfilePhotoURL = HandleNullString(profilePhotoURL)

	if dateOfBirth.Valid {
		user.DateOfBirth.Year = int32(dateOfBirth.Time.Year())
		user.DateOfBirth.Month = int32(dateOfBirth.Time.Month())
		user.DateOfBirth.Day = int32(dateOfBirth.Time.Day())
	}

	return user, nil
}
