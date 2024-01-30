package errors

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
)

// middleware
const (
	AuthorizationTokenNotProvided = "Authorization token not provided"
	RoleNotFoundInTokenClaims     = "Role not found in token claims"
	InsufficientPermission        = "Insufficient permissions"
	TokenClaimsNotFound           = "Token claims not found"
)
