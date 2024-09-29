package user

import (
	"context"
	"golang_lessons/pkg/logging"
)

type Service struct {
	logger  *logging.Logger
	storage Storage
}

func Create(ctx context.Context, dto CreateUserDTO) (u *User, err error) {

	return
}
