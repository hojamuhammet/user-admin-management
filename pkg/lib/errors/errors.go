package errors

import "errors"

// auth
const (
	InvalidRequestFormat    = "Invalid request format"
	AdminNotFound           = "Admin not found"
	InvalidCredentials      = "Invalid credentials"
	RefreshTokenNotProvided = "Refresh token not provided"
	InvalidRefreshToken     = "Invalid refresh token"
	InvalidURLParameters    = "Invalid URL parameters"
)

// user & admin
const (
	InternalServerError      = "Internal server error"
	InvalidID                = "Invalid ID"
	InvalidRequestBody       = "Invalid request body"
	InvalidPhoneNumberFormat = "Invalid phone number format"
	SearchQueryRequired      = "Search query is required"
	UserNotFound             = "User not found"
	PhoneNumberAlreadyInUse  = "Phone number already in use"
	EmailAlreadyInUse        = "Email already in use"
)

var ErrInternalServerError = errors.New("internal server error")
var ErrDatabaseError = errors.New("database error")

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrAdminNotFound        = errors.New("admin not found")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrAdminAlreadyExists   = errors.New("admin already exists")
	ErrAdminCannotBeDeleted = errors.New("super admin cannot be deleted")
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrPhoneNumberInUse   = errors.New("phone number already in use")
	ErrEmailInUse         = errors.New("email already in use")
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
)

// middleware
const (
	AuthorizationTokenNotProvided = "Authorization token not provided"
	RoleNotFoundInTokenClaims     = "Role not found in token claims"
	InsufficientPermission        = "Insufficient permissions"
	TokenClaimsNotFound           = "Token claims not found"
)
