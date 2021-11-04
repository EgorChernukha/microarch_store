package domain

type UserID int

type User interface {
	ID() UserID
	Login() string
	Firstname() string
	Lastname() string
	Email() string
	Phone() string
	Update(firstname, lastname, email, phone string)
}

type user struct {
	id        UserID
	login     string
	firstname string
	lastname  string
	email     string
	phone     string
}

func (u *user) ID() UserID {
	return u.id
}

func (u *user) Login() string {
	return u.login
}

func (u *user) Firstname() string {
	return u.firstname
}

func (u *user) Lastname() string {
	return u.lastname
}

func (u *user) Email() string {
	return u.email
}

func (u *user) Phone() string {
	return u.phone
}

func (u *user) Update(firstname, lastname, email, phone string) {
	u.firstname = firstname
	u.lastname = lastname
	u.email = email
	u.phone = phone
}

func NewUser(login, firstname, lastname, email, phone string) User {
	return &user{
		login:     login,
		firstname: firstname,
		lastname:  lastname,
		email:     email,
		phone:     phone,
	}
}
