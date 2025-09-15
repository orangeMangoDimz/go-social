package protocol

import "net/http"

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	WriteJSONError(w, http.StatusInternalServerError, "Something went wrong")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func ConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	WriteJSONError(w, http.StatusConflict, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	WriteJSONError(w, http.StatusNotFound, "not found")
}

func UnauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	WriteJSONError(w, http.StatusForbidden, "forbidden")
}

func UnauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	w.Header().Set("Retry-After", retryAfter)
	WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
