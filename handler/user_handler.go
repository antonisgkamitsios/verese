package handler

import (
	"net/http"

	"github.com/antonisgkamitsios/verese/internal/models"
	"github.com/antonisgkamitsios/verese/views/pages"
)

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {

	h.app.Render(w, r, http.StatusOK, pages.Register(models.User{}))
}

func (h *Handlers) RegisterPost(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name           string `form:"name"`
		Email          string `form:"email"`
		Password       string `form:"password"`
		PasswordRepeat string `form:"passwordRepeat"`
	}

	var form input

	err := h.app.DecodeForm(r, &form)
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}

	var user models.User

	user.Name = form.Name
	user.Email = form.Email
	user.Password.Set(form.Password)
	if form.PasswordRepeat == "" {
		user.AddError("passwordRepeat", "must be provided")
	}
	if form.Password != form.PasswordRepeat {
		user.AddError("passwordRepeat", "the passwords must match")
	}

	user.Validate()
	if !user.Valid() {
		h.app.Render(w, r, http.StatusUnprocessableEntity, pages.RegisterForm(user))
		return
	}

	h.app.Redirect(w, r, "/")
}
