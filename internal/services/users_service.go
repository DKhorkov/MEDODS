package services

import (
	"github.com/DKhorkov/medods/internal/interfaces"
)

type CommonUsersService struct {
	UsersRepository interfaces.UsersRepository
}

func (service *CommonUsersService) GetUserEmail(guid string) (string, error) {
	return service.UsersRepository.GetUserEmail(guid)
}
