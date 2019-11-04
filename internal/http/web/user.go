package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BoilerMake/bm-app/internal/models"
	"github.com/BoilerMake/bm-app/pkg/flash"

	"github.com/go-chi/chi"
)

// getSignup renders the signup template.
func (h *Handler) getSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Signup")

		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
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

	confirmMessage :=
		`Hey %s,

Thanks for creating a BoilerMake.org account! We're excited that you're interested in attending our hackathon.  Click the link below to confirm your email.

Be sure to check 'My Application' to ensure you've applied to BoilerMake VII.  This email only confirms that you've created a BoilerMake.org account – the application is separate.

%s

BoilerMake
Hack Your Own Adventure`

	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		id, confirmationCode, err := h.UserService.Signup(&u)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}
		u.ID = id

		// Build confirmation email
		to := u.Email
		subject := "Confirm your email"
		link := domain + "/activate/" + confirmationCode

		err = h.Mailer.Send(to, subject, fmt.Sprintf(confirmMessage, u.FirstName, link))
		if err != nil {
			h.Error(w, r, err, "", to, link)
			return
		}

		session := h.getSession(r)

		u.SetSession(session)
		err = session.Save(r, w)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		// Redirect to application if signup was successful
		http.Redirect(w, r, "/apply", http.StatusSeeOther)
	}
}

// getActivate activates the account that corresponds to the activation code
// if there is such an account.
func (h *Handler) getActivate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		u, err := h.UserService.GetByCode(code)
		if err != nil {
			h.Error(w, r, err, "/", code)
			return
		}

		u.IsActive = true
		u.ConfirmationCode = ""
		err = h.UserService.Update(u)
		if err != nil {
			h.Error(w, r, err, "/")
			return
		}

		// Redirect to application if activation was successful
		http.Redirect(w, r, "/apply", http.StatusSeeOther)
	}
}

// getForgotPassword renders the forgot password page.
func (h *Handler) getForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Forgot Password")

		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
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

	resetMessage :=
		`Hey there,

We got a request to reset your BoilerMake account's password, If you made this request then please click the button below.  If you didn't then feel free to ignore this email.

%s

BoilerMake
Hack Your Own Adventure`

	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		u.FromFormData(r)

		token, err := h.UserService.GetPasswordReset(u.Email)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		to := u.Email
		subject := "Password Reset"
		link := domain + "/reset/" + token

		err = h.Mailer.Send(to, subject, fmt.Sprintf(resetMessage, link))
		if err != nil {
			h.Error(w, r, err, "", to, link)
			return
		}

		session := h.getSession(r)

		// Show flash that they should be sent a reset email
		session.AddFlash(flash.Flash{
			Type:    flash.Info,
			Message: "We've emailed you instructions on how to reset your password.  If you don't see it in the next few mintues be sure to check your spam folder.",
		})
		session.Save(r, w)

		// Redirect to homepage if activation was successful
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getResetPassword renders the reset password template.
func (h *Handler) getResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Reset Password")

		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
			return
		}

		h.Templates.RenderTemplate(w, "reset", p)
	}
}

// getResetPasswordWithToken renders the reset password template with the token filled in.
func (h *Handler) getResetPasswordWithToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Reset Password")

		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
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
			h.Error(w, r, err, "")
			return
		}

		session := h.getSession(r)

		// Show flash that everything went well
		session.AddFlash(flash.Flash{
			Type:    flash.Success,
			Message: "Your password has been reset",
		})
		session.Save(r, w)

		// Redirect to login if reset was successful
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// getLogin renders the login template.
func (h *Handler) getLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := h.NewPage(w, r, "BoilerMake - Login")

		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
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
			h.Error(w, r, err, "")
			return
		}

		session := h.getSession(r)

		u.SetSession(session)
		err = session.Save(r, w)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		// Redirect to application if login was successful
		http.Redirect(w, r, "/apply", http.StatusSeeOther)
	}
}

// getLogout renders the login template.
func (h *Handler) getLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := h.getSession(r)

		// This expires the token
		session.Options.MaxAge = -1

		err := session.Save(r, w)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// getAccount shows a user their account.
func (h *Handler) getAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := h.getSession(r)

		email, ok := session.Values["EMAIL"].(string)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		u, err := h.UserService.GetByEmail(email)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		p, ok := h.NewPage(w, r, "BoilerMake - Account")
		if !ok {
			h.Error(w, r, errors.New("creating page failed"), "")
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
		session := h.getSession(r)

		email, ok := session.Values["EMAIL"].(string)
		if !ok {
			h.Error(w, r, errors.New("invalid session value"), "")
			return
		}

		u, err := h.UserService.GetByEmail(email)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		u.FirstName = r.FormValue("first-name")
		u.LastName = r.FormValue("last-name")

		err = h.UserService.Update(u)
		if err != nil {
			h.Error(w, r, err, "")
			return
		}

		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}
}
