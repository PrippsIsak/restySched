package handler

import (
	"net/http"

	"github.com/isak/restySched/web/templates"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(w http.ResponseWriter, r *http.Request) {
	templates.Home().Render(r.Context(), w)
}
