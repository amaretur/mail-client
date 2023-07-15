package handler

import (
	"strconv"
	"net/http"
)

func getOffsetAndLimit(r *http.Request) (int, int) {

	query := r.URL.Query()

	offsetStr := query.Get("offset")
	limitStr := query.Get("limit")

	if offsetStr == "" && limitStr == "" {
		return 0, 0
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return 0, 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		return 0, 0
	}

	return offset, limit
}
