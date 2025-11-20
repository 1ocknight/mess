package domain

type Service interface {
	CodeService() CodeService
	TokenService() TokenService
}

type Domain struct{
	CodeService 
	TokenService
}

func NewDomain(codeSvc CodeService, tokenSvc TokenService) *Domain {
	return &Domain{
		CodeService:  codeSvc,
		TokenService: tokenSvc,
	}
}