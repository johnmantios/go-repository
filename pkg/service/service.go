package service

import (
	"johnmantios.com/go-repository/pkg/repo"
)

type GreetingUserService struct {
	repository repo.IRepository
}

func NewGreetingUserService(repository repo.IRepository) GreetingUserService {
	return GreetingUserService{
		repository: repository,
	}
}

func (g *GreetingUserService) Greet(username string) (*repo.User, error) {

	user, err := g.repository.GetAUser(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
