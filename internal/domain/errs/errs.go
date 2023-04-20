package errs

import "fmt"

type DomainError struct {
	Msg string
}

func (e *DomainError) Error() string {
	return e.Msg
}

type DomainNotFoundError DomainError

func (en *DomainNotFoundError) Error() string {
	return en.Msg
}

func NewDomainNotFoundError(uuid string) *DomainNotFoundError {
	return &DomainNotFoundError{
		Msg: fmt.Sprintf("record not found witn uuid = %s", uuid),
	}
}
