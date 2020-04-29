package main

import "dite.pro/rollout-status/pkg/status"

type Output struct {
	Success bool        `json:"success"`
	Error   ErrorOutput `json:"error"`
}

type ErrorType string

const (
	ErrorTypeProgram ErrorType = "program"
	ErrorTypeRollout ErrorType = "rollout"
)

type ErrorOutput struct {
	Code    status.Failure `json:"code"`
	Message string         `json:"message"`
	Type    ErrorType      `json:"type"`
}

func ErrorOutputFrom(err error) ErrorOutput {
	if re, ok := err.(status.RolloutError); ok {
		return ErrorOutput{
			Code:    re.Failure,
			Message: re.Message,
			Type:    ErrorTypeRollout,
		}
	}
	return ErrorOutput{
		Message: err.Error(),
		Type:    ErrorTypeProgram,
	}
}

func RolloutOutput(rollout status.RolloutStatus) Output {
	out := Output{
		Success: rollout.Error == nil,
		Error:   ErrorOutputFrom(rollout.Error),
	}

	return out
}
