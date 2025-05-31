package user

import (
	"time"
)

type UserService interface {
	LoginUser(email string, password string) (string, time.Time, error)
}

type ServiceDependencies struct {
	UserRepository UserRepository
	GenerateToken  func(map[string]any) (string, time.Time, error)
}

type userService struct {
	Dependencies ServiceDependencies
}

func NewUserService(dependencies ServiceDependencies) userService {
	return userService{
		Dependencies: dependencies,
	}
}

func (service *userService) LoginUser(email string, password string) (string, time.Time, error) {
	user, err := service.Dependencies.UserRepository.GetByEmail(email)
	if err != nil || user.Password != password {
		return "", time.Time{}, ErrLoginFailed
	}

	token, expires_at, err := service.Dependencies.GenerateToken(map[string]any{"user_id": user.Id})
	if err != nil {
		return "", time.Time{}, ErrTokenGenerationFailed
	}

	return token, expires_at, nil
}
