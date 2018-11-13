package app

import (
	"github.com/go-chi/chi"
)

// routes generates a router and assigns it to the Server's handler.
// It will overwrite any handler that may already exist.
func (s *Server) routes() {
	r := chi.NewRouter()

	r.Get("/", s.getRoot())
	/*
		r.Get("/faq", s.getFAQ)
		r.Get("/login", s.getLogin)
		r.Post("/login", s.postLogin)
		r.Get("/register", s.getRegister)
		r.Post("/register", s.postRegister)
		// TODO GitHub auth

		r.Route("/api", func(r chi.Router) {
			r.Get("/announcements", s.getAnnouncements)
			r.Route("/user", func(r chi.Router) {
				r.Post("/login", s.login)
				r.Post("/register", s.register)
				r.Get("/confirm/{code}", s.confirmCode)
			})
		})

		r.Route("/exec", func(r chi.Router) {
			r.Use(s.execOnly())

			r.Get("/", s.execRoot)

			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.execGetUsers)
				r.Post("/search", s.execSearchUsers)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", s.execGetUser)
					r.Post("/{id}/passwordreset", s.execResetPassword)
					r.Post("/{id}/checkin", s.execCheckinUser)
				})
			})

			// TODO should execs be able to edit apps? I think so
			r.Routes("/applications", func(r chi.Router) {
				r.Get("/", s.execGetApplications)
				r.Get("/{id}", s.execGetApplication)
			})
		})
	*/

	s.Handler = r
}

/*
Current routes based on front end

/
/sponsors
/hackers
/about
/live
/contact
/faq

/register
/login
/reset/{reset_token}
/confirm/{confirm_token}

/exec
/exec/checkin
/exec/users
/exec/users/{userid}
/exec/applications
/exec/applications/{appid}

Routes I think we'll need based on backend + frontend

*/
