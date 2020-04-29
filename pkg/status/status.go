package status

import (
	"fmt"
)

type RolloutStatus struct {
	Continue bool  `json:"continue"`
	Error    error `json:"error"`
}

// Rollout is not completed and it will not succeed with a high certainty.
// Also returned for program issues not directly relevant to the rollout.
func RolloutFatal(err error) RolloutStatus {
	return RolloutStatus{
		Continue: false,
		Error:    err,
	}
}

// Rollout is not completed but there is a chance it will eventually succeed.
// Returned for valid flow (ContainerCreating -> init containers -> not ready -> ready)
// as well as for error states under a deadline.
func RolloutErrorProgressing(err error) RolloutStatus {
	return RolloutStatus{
		Continue: true,
		Error:    err,
	}
}

// Rollout completed successfully
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
