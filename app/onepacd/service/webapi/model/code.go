package model

import (
	"errors"
	"fmt"
)

const (
	Code_Success = 200

	code_Common_Start  = 1000
	Code_InternalError = code_Common_Start + 0
	Code_InvalidParams = code_Common_Start + 1
	Code_NotFound      = code_Common_Start + 2
	Code_Unauthorized  = code_Common_Start + 3
	Code_Forbidden     = code_Common_Start + 4
	Code_TooMany       = code_Common_Start + 5
	Code_Busy          = code_Common_Start + 6
	Code_UnknownError  = code_Common_Start + 7
	Code_DatabaseError = code_Common_Start + 8

	// Assets

	Code_Assets_EpochRunError = 2000
)

var (
	errorSets = []map[int]error{
		commonErrors,
	}

	unknownError = errors.New("unknown error")

	commonErrors = map[int]error{
		Code_Success:       errors.New("success"),
		Code_InternalError: errors.New("internal error"),
		Code_InvalidParams: errors.New("invalid param"),
		Code_NotFound:      errors.New("not found"),
		Code_Unauthorized:  errors.New("unauthorized"),
		Code_Forbidden:     errors.New("forbidden"),
		Code_TooMany:       errors.New("too many"),
		Code_Busy:          errors.New("busy"),
		Code_UnknownError:  unknownError,
		Code_DatabaseError: errors.New("database error"),
	}

	allCodeToErrors = map[int]error{}
	allErrorsToCode = map[error]int{}
)

func init() {
	for _, errs := range errorSets {
		for code, err := range errs {
			if _, ok := allCodeToErrors[code]; ok {
				panic(fmt.Sprintf("error code %d is already registered", code))
			}
			allCodeToErrors[code] = err
		}
	}

	for code, err := range allCodeToErrors {
		allErrorsToCode[err] = code
	}
}

func ErrorFromCode(code int) error {
	err, ok := allCodeToErrors[code]

	if !ok {
		return unknownError
	}

	return err
}

func CodeFromError(err error) int {
	code, ok := allErrorsToCode[err]

	if !ok {
		return Code_InternalError
	}

	return code
}
