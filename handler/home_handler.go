package handler

import (
	"net/http"

	"github.com/antonisgkamitsios/verese/views/pages"
)

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {

	h.app.Render(w, r, http.StatusOK, pages.Home())
}
