package handler

import (
	"fmt"
	"net/http"

	"github.com/antonisgkamitsios/verese/internal/reqctx"
	"github.com/antonisgkamitsios/verese/views/pages"
)

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {

	user := reqctx.User(r.Context())
	fmt.Printf("user: %+v\n", user)

	h.app.Render(w, r, http.StatusOK, pages.Home())
}
