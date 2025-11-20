package identifier

type TokenInfo interface {
	GetAccessToken() string
	GetRefreshToken() string
	GetExpiresIn() int
	GetRefreshExpiresIn() int
	GetTokenType() string
}

type Service interface {
	ExchangeCode(code string, redirectURL string) (TokenInfo, error)
	Refresh(refreshToken string) (TokenInfo, error)
}
