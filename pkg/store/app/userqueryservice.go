package app

type UserData struct {
	ID        int
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Phone     string
}

type UserQueryService interface {
	FindUser(id int) (UserData, error)
}
