package web

import (
	"net/http"

	"github.com/BoilerMake/new-backend/internal/http/middleware"
	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
)

// getSignup renders the signup template.
func (h *Handler) getSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - Sign Up", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "signup", p)
	}
}

// postSignup tries to signup a user from a post request.
func (h *Handler) postSignup() http.HandlerFunc {
	domain := mustGetEnv("DOMAIN")
	mode := mustGetEnv("ENV_MODE")
	if mode == "development" {
		domain = "http://" + domain + ":" + mustGetEnv("PORT")
	} else {
		domain = "https://" + domain
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO check if login is valid (i.e. account exists), if so log them in
		var u models.User
		u.FromFormData(r)

		id, confirmationCode, err := h.UserService.Signup(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.ID = id

		// Build confirmation email
		to := u.Email
		subject := "Confirm your email"
		link := domain + "/activate/" + confirmationCode
		data := map[string]interface{}{
			"Name":        u.FirstName,
			"ConfirmLink": link,
		}

		err = h.Mailer.SendTemplate(to, subject, "email confirm", data)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			// TODO error handling, this state should never be reached
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		u.SetSession(session)
		err = session.Save(r, w)
		if err != nil {
			// TODO Error Handling, this state should never be reached
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to homepage if signup was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getActivate activates the account that corresponds to the activation code
// if there is such an account.
func (h *Handler) getActivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		u, err := h.UserService.GetByCode(code)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u.IsActive = true
		u.ConfirmationCode = ""
		err = h.UserService.Update(u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// TODO once session tokens are updated this should show a success flash
		// Redirect to homepage if activation was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getForgotPassword renders the forgot password page.
func (h *Handler) getForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - Forgot Password", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "forgot", p)
	}
}

// postForgotPassword sends the password reset email.
func (h *Handler) postForgotPassword() http.HandlerFunc {
	domain := mustGetEnv("DOMAIN")
	mode := mustGetEnv("ENV_MODE")
	if mode == "development" {
		domain = "http://" + domain + ":" + mustGetEnv("PORT")
	} else {
		domain = "https://" + domain
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		token, err := h.UserService.GetPasswordReset(u.Email)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		to := u.Email
		subject := "Password Reset"
		link := domain + "/reset/" + token
		data := map[string]interface{}{
			"Name":      u.FirstName,
			"ResetLink": link,
		}

		err = h.Mailer.SendTemplate(to, subject, "email reset", data)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO once session tokens are updated this should show a success flash
		// Redirect to homepage if activation was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getResetPassword renders the reset password template.
func (h *Handler) getResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - Reset Password", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "reset", p)
	}
}

// getResetPasswordWithToken renders the reset password template with the token filled in.
func (h *Handler) getResetPasswordWithToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - Reset Password", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		p.Data = map[string]interface{}{
			"Token": chi.URLParam(r, "token"),
		}

		h.Templates.RenderTemplate(w, "reset", p)
	}
}

// postResetPassword resets the password with a valid token
func (h *Handler) postResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var passwordResetInfo models.PasswordResetPayload
		passwordResetInfo.UserToken = r.FormValue("token")
		passwordResetInfo.NewPassword = r.FormValue("new-password")

		err := h.UserService.ResetPassword(passwordResetInfo.UserToken, passwordResetInfo.NewPassword)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO once session tokens are updated this should show a success flash
		// Redirect to homepage if activation was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getLogin renders the login template.
func (h *Handler) getLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := NewPage("BoilerMake - Login", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		h.Templates.RenderTemplate(w, "login", p)
	}
}

// postLogin tries to log in a user.
func (h *Handler) postLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		err := h.UserService.Login(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			// TODO error handling, this state should never be reached
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		u.SetSession(session)
		err = session.Save(r, w)
		if err != nil {
			// TODO Error Handling, this state should never be reached
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to homepage if login was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getLogout renders the login template.
func (h *Handler) getLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		// This expires the token
		session.Options.MaxAge = -1

		err := session.Save(r, w)
		if err != nil {
			// TODO Error Handling, this state should never be reached
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getAccount shows a user their account.
func (h *Handler) getAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		email, ok := session.Values["EMAIL"].(string)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "invalid session value", http.StatusInternalServerError)
			return
		}

		u, err := h.UserService.GetByEmail(email)
		if err != nil {
			// TODO error handling
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		p, ok := NewPage("BoilerMake - Account", r)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "creating page failed", http.StatusInternalServerError)
			return
		}

		p.FormRefill = map[string]interface{}{
			"FirstName": u.FirstName,
			"LastName":  u.LastName,
			"Email":     u.Email,
		}

		h.Templates.RenderTemplate(w, "account", p)
	}
}

// postAccount updates a user's account.
func (h *Handler) postAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(middleware.SessionCtxKey).(*sessions.Session)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "getting session failed", http.StatusInternalServerError)
			return
		}

		email, ok := session.Values["EMAIL"].(string)
		if !ok {
			// TODO Error Handling, this state should never be reached
			http.Error(w, "invalid session value", http.StatusInternalServerError)
			return
		}

		u, err := h.UserService.GetByEmail(email)
		if err != nil {
			// TODO error handling
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		u.FirstName = r.FormValue("first-name")
		u.LastName = r.FormValue("last-name")

		err = h.UserService.Update(u)
		if err != nil {
			// TODO error handling
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}
}
