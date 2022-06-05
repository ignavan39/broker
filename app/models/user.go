package models

type User struct {
	Id string `json:"id"`
	Password string `json:"password"`
	Email string `json:"email"`
	LastName string `json:"lastName"`
	FirstName string `json:"firstName"`
}