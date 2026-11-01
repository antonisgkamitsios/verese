package handler

import (
	"net/http"

	"github.com/antonisgkamitsios/verese/views"
)

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	views.Home().Render(r.Context(), w)
}
