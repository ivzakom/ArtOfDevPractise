package user

import (
	"artOfDevPractise/pkg/logging"
	"context"
)

type Service struct {
	logger  *logging.Logger
	storage Storage
}

func Create(ctx context.Context, dto CreateUserDTO) (u *User, err error) {

	return
}
