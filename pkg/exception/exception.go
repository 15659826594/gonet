package exception

import (
	"errors"
	"fmt"
)

var (
	Exception             = errors.New("base exception")
	RuntimeException      = fmt.Errorf("%w: runtime exception", Exception)
	HttpResponseException = fmt.Errorf("%w: http response exception", RuntimeException)
)

type Throwable interface {
	Exception() error
}
