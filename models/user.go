package models

type User struct {
	Base
	Email string

	Name     string
	Provider *string
	Avi      string
	Sex      string
	Username string
}
