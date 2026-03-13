package myerrors

import (
	"errors"
	"fmt"
	"shb/internal/models"
)

var (
	ErrGeneral  = errors.New("service temporarily unavailable")
	ErrNotFound = errors.New("not found")
)

func InternalError() models.Response {
	return models.Response{
		Message: ErrGeneral.Error(),
	}
}

type BadRequestErr struct {
	Message string `json:"message"`
	err     error
}

func (e BadRequestErr) Error() string {
	return fmt.Sprintf("message %s: %v", e.Message, e.err)
}

type ForbiddenErr struct {
	Message string `json:"message"`
	err     error
}

func (e ForbiddenErr) Error() string {
	return fmt.Sprintf("message %s: %v", e.Message, e.err)
}

type UnprocessableErr struct {
	Message string `json:"message"`
	err     error
}

func (e UnprocessableErr) Error() string {
	return fmt.Sprintf("message %s: %v", e.Message, e.err)
}

type UnauthorizedErr struct {
	Message string `json:"message"`
	err     error
}

func (e UnauthorizedErr) Error() string {
	return fmt.Sprintf("message %s: %v", e.Message, e.err)
}

type TooManyRequestsErr struct {
	Message string `json:"message"`
	err     error
}

func (e TooManyRequestsErr) Error() string {
	return fmt.Sprintf("message %s: %v", e.Message, e.err)
}

type ConflictErr struct {
	Message string `json:"message"`
	err     error
}

func (e ConflictErr) Error() string {
	return fmt.Sprintf("message %s: %v", e.Message, e.err)
}

func NewBadRequestErr(message string) error {
	return BadRequestErr{Message: message}
}

func NewForbiddenErr(message string) error {
	return ForbiddenErr{Message: message}
}

func NewUnprocessableErr(message string) error {
	return UnprocessableErr{Message: message}
}

func NewUnauthorizedErr(message string) error {
	return UnauthorizedErr{Message: message}
}

func NewTooManyRequestsErr(message string) error {
	return TooManyRequestsErr{Message: message}
}

func NewConflictErr(message string) error {
	return ConflictErr{Message: message}
}
