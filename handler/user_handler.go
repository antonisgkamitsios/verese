package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/antonisgkamitsios/verese/internal/models"
	"github.com/antonisgkamitsios/verese/views/mail"
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

	err = h.app.Models.Users.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			user.AddError("email", "this email is already in use")
			h.app.Render(w, r, http.StatusUnprocessableEntity, pages.RegisterForm(user))
			return
		default:
			h.app.ServerError(w, r, err)
			return
		}
	}

	h.app.Background(func() {
		err = h.app.Mailer.Send(user.Email, "Welcome lover", mail.Hello(user))
		if err != nil {
			h.app.Logger.Error(fmt.Sprintf("Error sending email %v", err))
		}
	})

	h.app.Redirect(w, r, "/")
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {

	h.app.Render(w, r, http.StatusOK, pages.Login(models.LoginForm{}))
}
func (h *Handlers) LoginPost(w http.ResponseWriter, r *http.Request) {

	var form models.LoginForm
	err := h.app.DecodeForm(r, &form)
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}
	form.Validate()
	if !form.Valid() {
		h.app.Render(w, r, http.StatusUnprocessableEntity, pages.LoginForm(form))
		return
	}

	user, err := h.app.Models.Users.GetByEmail(form.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			form.AddNonFieldError("invalid credentials")
			h.app.Render(w, r, http.StatusUnprocessableEntity, pages.LoginForm(form))
			return
		default:
			h.app.ServerError(w, r, err)
			return

		}
	}

	passwordMatches, err := user.Password.Matches(form.Password)
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}
	if !passwordMatches {
		form.AddNonFieldError("invalid credentials")
		h.app.Render(w, r, http.StatusUnprocessableEntity, pages.LoginForm(form))
		return
	}

	err = h.app.Session.RenewToken(r.Context())
	if err != nil {
		h.app.ServerError(w, r, err)
		return

	}
	h.app.Session.Put(r.Context(), "user_id", user.ID)

	h.app.Redirect(w, r, "/")
}

func (h *Handlers) LogoutPost(w http.ResponseWriter, r *http.Request) {
	err := h.app.Session.RenewToken(r.Context())
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}

	h.app.Session.Remove(r.Context(), "user_id")
	// todo when remember me is implemented here we will delete token + the cookie

	h.app.Redirect(w, r, "/")
}
