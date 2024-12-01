package repo

type IRepository interface {
	GetAUser(string) (*User, error)
}

type User struct {
	Name string `json:"name"`
}
