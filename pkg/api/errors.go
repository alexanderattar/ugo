package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse type for handling HTTP errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render an API error response
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// Error400 renders an Invalid Request response
func Error400(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid Request",
		ErrorText:      err.Error(),
	}
}

// Error500 renders an Internal Server Error response
func Error500(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}

// Error404 renders an Not Found response
var Error404 = &ErrResponse{HTTPStatusCode: 404, StatusText: "Not Found"}

// CustomError renders an Invalid Request response
func CustomError(err string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 400,
		StatusText:     "Invalid Request",
		ErrorText:      err,
	}
}
