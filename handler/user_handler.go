package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/antonisgkamitsios/verese/internal/models"
	"github.com/antonisgkamitsios/verese/internal/reqctx"
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

	token, err := h.app.Models.Tokens.New(user.ID, 24*time.Hour, models.ScopeActivation)
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}

	h.app.Background(func() {
		err = h.app.Mailer.Send(user.Email, "Welcome lover", mail.Hello(user, *token, h.app.Config.BaseUrl))
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

	if !user.Activated {
		form.AddNonFieldError("you must activate your account first! Visit your email for instructions")
		h.app.Render(w, r, http.StatusUnprocessableEntity, pages.LoginForm(form))
		return
	}

	err = h.app.Session.RenewToken(r.Context())
	if err != nil {
		h.app.ServerError(w, r, err)
		return

	}
	h.app.Session.Put(r.Context(), "user_id", user.ID)

	if form.RememberMe {
		rememberMeTok, err := h.app.Models.Tokens.New(user.ID, 365*24*time.Hour, models.ScopeRememberMe)
		if err != nil {
			h.app.ServerError(w, r, err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "remember_me",
			Value:    rememberMeTok.PlainText,
			Path:     "/",
			MaxAge:   int((365 * 24 * time.Hour) / time.Second),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})

	} else {
		err = h.app.Models.Tokens.DeleteAllForUser(user.ID, models.ScopeRememberMe)
		if err != nil {
			h.app.ServerError(w, r, err)
			return
		}
		cookie, _ := r.Cookie("remember_me")
		if cookie != nil {
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		}
	}

	h.app.Redirect(w, r, "/")
}

func (h *Handlers) LogoutPost(w http.ResponseWriter, r *http.Request) {
	err := h.app.Session.RenewToken(r.Context())
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}

	user := reqctx.CurrentUser(r.Context())
	if user.IsAnonymous() {
		h.app.Redirect(w, r, "/")
		return
	}

	h.app.Session.Remove(r.Context(), "user_id")
	err = h.app.Models.Tokens.DeleteAllForUser(user.ID, models.ScopeRememberMe)
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}
	cookie, _ := r.Cookie("remember_me")
	if cookie != nil {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}

	h.app.Redirect(w, r, "/")
}

func (h *Handlers) UserActivate(w http.ResponseWriter, r *http.Request) {
	id, err := h.app.ParseIDParam(r)
	if err != nil {
		h.app.ClientError(w, r, http.StatusNotFound, "the resource was not found")
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		h.app.ClientError(w, r, http.StatusNotFound, "the resource was not found")
		return
	}

	user, err := h.app.Models.Users.FindByToken(models.ScopeActivation, token)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			h.app.ClientError(w, r, http.StatusNotFound, "the resource was not found")
			return
		default:
			h.app.ServerError(w, r, err)
			return
		}
	}

	if user.ID != id {
		h.app.ClientError(w, r, http.StatusNotFound, "the resource was not found")
		return
	}

	user.Activated = true
	err = h.app.Models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			h.app.ClientError(w, r, http.StatusNotFound, "there was an error please try again")
		default:
			h.app.ServerError(w, r, err)
		}
		return
	}

	err = h.app.Models.Tokens.DeleteAllForUser(user.ID, models.ScopeActivation)
	if err != nil {
		h.app.ServerError(w, r, err)
		return
	}

	fmt.Fprintf(w, "User  is: %+v\n", user)
}
