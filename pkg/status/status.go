package status

import (
	"bytes"
	"fmt"
	"strings"
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

type ErrorGroup struct {
	Errors []error `json:"errors"`
}

func (group ErrorGroup) Error() string {
	var out bytes.Buffer
	for _, err := range group.Errors {
		fmt.Fprintf(&out, "%v\n", err.Error())
	}
	return strings.TrimSpace(out.String())
}

func (group *ErrorGroup) Add(err error) {
	group.Errors = append(group.Errors, err)
}

func (group ErrorGroup) Empty() bool {
	return len(group.Errors) == 0
}
