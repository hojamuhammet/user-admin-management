package service

type AdminAuthServiceInterface interface {
	LoginAdmin(username, password string) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
	LogoutAdmin(refreshToken string) error
}
