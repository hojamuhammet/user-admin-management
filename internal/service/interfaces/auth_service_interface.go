package service

type AuthService interface {
	LoginAdmin(username, password string) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
	LogoutAdmin(refreshToken string) error
}
