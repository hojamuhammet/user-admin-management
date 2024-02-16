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

func IsValidPhoneNumber(phoneNumber string) bool {
	const validPrefix = "+993"
	return len(phoneNumber) == 12 && strings.HasPrefix(phoneNumber, validPrefix)
}

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
		slog.Error("Error scanning user row: %v", Err(err))
		return domain.CommonUserResponse{}, err
	}

	return user, nil
}
