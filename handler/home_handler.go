package handler

import (
	"fmt"
	"net/http"

	"github.com/antonisgkamitsios/verese/internal/csrf"
	"github.com/antonisgkamitsios/verese/views/pages"
)

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {

	val := r.Context().Value(csrf.CSRFContextKey).(string)
	fmt.Printf("CSRF: %s\n", val)
	h.app.Render(w, r, http.StatusOK, pages.Home())
}
