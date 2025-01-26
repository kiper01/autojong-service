package handler

import (
	"net/http"
	"strings"
)

func errorCode(err error) int {

	if err == nil {
		return http.StatusOK
	}

	errMsg := err.Error()
	var httpStatus int

	switch {
	case strings.Contains(errMsg, "no rows"), strings.Contains(errMsg, "not found"):
		httpStatus = http.StatusNotFound
	case strings.Contains(errMsg, "already exists"), strings.Contains(errMsg, "duplicate key value violates unique constraint"):
		httpStatus = http.StatusConflict
	case strings.Contains(errMsg, "foreign key constraint"):
		httpStatus = http.StatusForbidden
	case strings.Contains(errMsg, "could not serialize access due to concurrent update"):
		httpStatus = http.StatusConflict
	case strings.Contains(errMsg, "unsupported command type"):
		httpStatus = http.StatusNotImplemented
	case strings.Contains(errMsg, "update statements must have at least one Set clause"):
		httpStatus = http.StatusBadRequest
	case strings.Contains(errMsg, "invalid"):
		httpStatus = http.StatusBadRequest
	case strings.Contains(errMsg, "conflict"):
		httpStatus = http.StatusConflict
	default:
		httpStatus = http.StatusInternalServerError
	}

	return httpStatus
}
