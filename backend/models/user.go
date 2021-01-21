package models

type User interface {
	GetId() string
	GetName() string
}

type UserRepository interface {
	AddUser(user User) error
	RemoveUser(user User)
	FindUserById(ID string) User
	GetAllUsers() []User
	GetAllRegisteredUsers() []User
}
