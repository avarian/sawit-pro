// This file contains types that are used in the repository layer.
package repository

import "time"

type CreateUserInput struct {
	FullName    string
	PhoneNumber string
	Password    string
}

type UpdateUserInput struct {
	Id          int
	FullName    string
	PhoneNumber string
}

type User struct {
	Id          int
	FullName    string
	PhoneNumber string
	Password    string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}
