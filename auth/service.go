package auth

type Service interface {
	Authorize(user string) bool
}

type alwaysTrueAuthService struct {
}

func NewAlwaysTrustEveryoneAuthService() Service {
	return &alwaysTrueAuthService{}
}

func (alwaysTrueAuthService) Authorize(user string) bool {
	return true
}
