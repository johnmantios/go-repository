package service

type UserGetter interface {
	GetAUser(string) (*User, error)
}

type User struct {
	Name string `json:"name"`
}

type GreetingUserService struct {
	userGetter UserGetter
}

func NewGreetingUserService(userGetter UserGetter) GreetingUserService {
	return GreetingUserService{
		userGetter: userGetter,
	}
}

func (g *GreetingUserService) Greet(username string) (*User, error) {

	user, err := g.userGetter.GetAUser(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
