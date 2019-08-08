package errors

import (
	"errors"
	"fmt"
	"io"

	"github.com/visionom/v-gae/log"
	"gopkg.in/mgo.v2/bson"
)

var DEBUG bool

type VError struct {
	types string
	cause error
	msg   string
	*stack
}

type stackTracer interface {
	StackTrace() StackTrace
}

func (r *VError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", r.Cause())
			r.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%+v\n", r.Cause())
		io.WriteString(s, r.Error())
	case 'q':
		fmt.Fprintf(s, "%q", r.Error())
		io.WriteString(s, r.Error())
	}
}

func (r *VError) Cause() error { return r.cause }

func (r *VError) Msg() string { return r.msg }

func (r VError) Error() string {
	return r.msg + ":" + r.cause.Error()
}

func New(types, message string) error {
	err := errors.New(message)
	return VError{
		types,
		err,
		message,
		callers(),
	}
}

func Newf(types, format string, args ...interface{}) error {
	err := errors.New(fmt.Sprintf(format, args...))
	return VError{
		types,
		err,
		fmt.Sprintf(format, args...),
		callers(),
	}
}

func Wrap(types string, err error, message string) error {
	if err == nil {
		return nil
	}
	return VError{
		types,
		err,
		message,
		callers(),
	}
}

func Wrapf(types string, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return VError{
		types,
		err,
		fmt.Sprintf(format, args...),
		callers(),
	}
}

func HandleError(ticketID string, err error) {
	if ticketID == "" {
		ticketID = bson.NewObjectId().Hex()
	}
	if err, ok := err.(stackTracer); DEBUG && ok {
		log.Ef("[TicketID: %v]: %+v", ticketID, err)
		for _, f := range err.StackTrace() {
			log.Ef("[TicketID: %v]: %+s:%d, func %n", ticketID, f, f, f)
		}
	} else {
		logStr := fmt.Sprintf("%+v", err)
		if len(logStr) > 200 {
			log.Ef("[TicketID: %v]: %s ... total %d chars", ticketID, logStr[0:200],
				len(logStr))
		} else {
			log.Ef("[TicketID: %v]: %s", ticketID, logStr)
		}
	}
}
