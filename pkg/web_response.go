package pkg

import (
	"net/http"
)

type ErrorResp struct {
	Error  string            `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}

// error uses web.WriteError to return an error message to the client.
func ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	WriteError(w, r, err)
}

// JsonResponse() is a helper for sending json responses.
func JsonResponse(w http.ResponseWriter, r *http.Request, status int, data any) {
	if err := WriteJSON(w, status, data, nil); err != nil {
		ErrorResponse(w, r, err)
	}
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, Errorf(ENOTFOUND, "Not found."))
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, errorMessage string) {
	JsonResponse(w, r, http.StatusBadRequest, &ErrorResp{Error: errorMessage})
}

func ValidationErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	resp := ErrorResp{Error: "Validation failed.", Errors: errors}
	JsonResponse(w, r, http.StatusUnprocessableEntity, resp)
}

func RateLimitedErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	resp := ErrorResp{Error: "Rate limited!", Errors: errors}
	JsonResponse(w, r, http.StatusTooManyRequests, resp)
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := ErrorCode(err), ErrorMessage(err)

	if err := WriteJSON(w, ErrorStatusCode(code), &ErrorResp{Error: message}, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
