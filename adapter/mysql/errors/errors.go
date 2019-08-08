package errors

import (
	"github.com/visionom/v-gae/errors"
)

const (
	types = "mysql"
)

func New(message string) error {
	return errors.New(types, message)
}

func Newf(format string, args ...interface{}) error {
	return errors.Newf(types, format, args...)
}

func Wrap(err error, message string) error {
	return errors.Wrap(types, err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(types, err, format, args...)
}
