package status

import (
	"fmt"
)

type RolloutStatus struct {
	Continue bool  `json:"continue"`
	Error    error `json:"error"`
}

func RolloutFatal(err error) RolloutStatus {
	return RolloutStatus{
		Continue: false,
		Error:    err,
	}
}

func RolloutErrorProgressing(err error) RolloutStatus {
	return RolloutStatus{
		Continue: true,
		Error:    err,
	}
}

func RolloutOk() RolloutStatus {
	return RolloutStatus{
		Continue: false,
		Error:    nil,
	}
}

type RolloutError struct {
	Message string `json:"message"`
}

func (re RolloutError) Error() string {
	return re.Message
}

func MakeRolloutErorr(format string, args ...interface{}) RolloutError {
	return RolloutError{
		Message: fmt.Sprintf(format, args...),
	}
}
