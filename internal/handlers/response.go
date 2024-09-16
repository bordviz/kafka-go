package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Detail any `json:"detail"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, detail interface{}) {
	render.Status(r, status)
	render.JSON(w, r, Response{Detail: detail})
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
