package filters

import (
	"errors"
	"fmt"
	"golang_lessons/internal/item"
)

type FiltersErr struct {
	Err     error
	ErrDesc string
}

func NewFilterError(err error, desc string) FiltersErr {
	return FiltersErr{Err: err, ErrDesc: desc}
}

func (e FiltersErr) Error() string {
	return e.ErrDesc
}

func (e FiltersErr) Unwrap() error {
	return e.Err
}

func FilterNotNull(item item.Item, quantity int) error {

	var err error

	if item.GetID() == 1 || item.GetID() == 2 {
		if quantity <= 0 {
			errDesc := fmt.Sprintf("The quantity of %s must be greater than zero.", item.GetName())
			err = errors.New(errDesc)
		}
	}
	return err
}
