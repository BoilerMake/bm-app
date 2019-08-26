package web

import (
	"net/http"

	"github.com/BoilerMake/new-backend/internal/models"

	"github.com/go-chi/chi"
)

// getSignup renders the signup template.
func (h *Handler) getSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Templates.RenderTemplate(w, "signup", nil)
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

	jwtIssuer := mustGetEnv("JWT_ISSUER")
	jwtSigningKey := []byte(mustGetEnv("JWT_SIGNING_KEY"))
	jwtCookie := mustGetEnv("JWT_COOKIE_NAME")

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

		jwt, err := u.GetJWT(jwtIssuer, jwtSigningKey)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO right now this is only valid on the domain it's sent from, if we do
		// subdomains (seems likely) then we'll need to change that.
		http.SetCookie(w, &http.Cookie{
			Name:       jwtCookie,
			Value:      jwt,
			Path:       "/",
			RawExpires: "0",
		})

		// Redirect to homepage if signup was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getActivate activates the account that corresponds to the activation code
// if there is such an account.
func (h *Handler) getActivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Temporarily ignoring claims returned from getClaimsFromCtx
		_, err := getClaimsFromCtx(r.Context())
		if err != nil {
			// TODO error handling
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

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
		h.Templates.RenderTemplate(w, "forgot", nil)
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
		h.Templates.RenderTemplate(w, "reset", nil)
	}
}

// getResetPasswordWithToken renders the reset password template with the token filled in.
func (h *Handler) getResetPasswordWithToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		h.Templates.RenderTemplate(w, "reset", token)
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
		h.Templates.RenderTemplate(w, "login", nil)
	}
}

// postLogin tries to log in a user.
func (h *Handler) postLogin() http.HandlerFunc {
	jwtIssuer := mustGetEnv("JWT_ISSUER")
	jwtSigningKey := []byte(mustGetEnv("JWT_SIGNING_KEY"))
	jwtCookie := mustGetEnv("JWT_COOKIE_NAME")

	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		err := h.UserService.Login(&u)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jwt, err := u.GetJWT(jwtIssuer, jwtSigningKey)
		if err != nil {
			// TODO error handling
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO Right now this is only valid on the domain it's sent from, if we do
		// subdomains (seems likely) then we'll need to change that.
		http.SetCookie(w, &http.Cookie{
			Name:       jwtCookie,
			Value:      jwt,
			Path:       "/",
			RawExpires: "0",
		})

		// Redirect to homepage if login was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getAccount shows a user their account.
func (h *Handler) getAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getClaimsFromCtx(r.Context())
		if err != nil {
			// TODO error handling
			// TODO once session tokens are updated this should show a need to login first flash
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		u, err := h.UserService.GetByEmail(claims["email"].(string))
		if err != nil {
			// TODO error handling
			// This can fail either because the DB is messed up or nothing is found
			// So be sure to deal with that
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		data := map[string]interface{}{
			"Email":     u.Email,
			"FirstName": u.FirstName,
			"LastName":  u.LastName,
			"Phone":     u.Phone,
		}

		h.Templates.RenderTemplate(w, "account", data)
	}
}
