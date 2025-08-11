package utils

import (
	"errors"

	"github.com/thecipherdev/goauth/mock"
	"github.com/thecipherdev/goauth/model"
)

func GetByUsername(username string) (*model.User, error) {
	for i := range mock.Users {
		if mock.Users[i].Username == username {
			return &mock.Users[i], nil
		}
	}
	return nil, errors.New("no user found")
}
