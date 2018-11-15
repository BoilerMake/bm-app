package app

import (
	"github.com/go-chi/chi"
)

// routes generates a router and assigns it to the Server's handler.
// It will overwrite any handler that may already exist in that server.
func (s *Server) routes() {
	r := chi.NewRouter()

	// NOTE Currently our api runs at api.boilermake.org.
	// But this has it running at boilermake.org/api
	// There's not much of a difference, but we may need to update frontend endpoints
	// NOTE we have versioning in our current backend, but I don't see much reason as to why
	// It's not exactly public, and would have little use people outside of BM
	// If we ever make an app, that could change but it seems not needed to me
	// Maybe that's too much thought into putting three more characters in URL, idk
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", s.getPing())
		/*
			// TODO add interest for exec and hacker

			// NOTE should we be versioning announcements?
			r.Get("/announcements", s.getAnnouncements)

			r.Route("/users", func(r chi.Router) {
				// TODO GitHub auth
				r.Post("/login", s.login)
				r.Post("/register", s.register)
				r.Get("/confirm/{code}", s.confirmCode)
			})

			r.Route("/exec", func(r chi.Router) {
				r.Use(s.execOnly())

				r.Get("/dashboard", s.execRoot)

				r.Route("/users", func(r chi.Router) {
					r.Get("/", s.execGetUsers)
					r.Post("/search", s.execSearchUsers)

					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", s.execGetUser)
						// TODO should execs be able to reset passwords for users?
						// I don't think so, if an exec account gets hack that could be bad
						// In our old backend, they were able to
						//r.Post("/{id}/passwordreset", s.execResetPassword)
						r.Post("/{id}/checkin", s.execCheckinUser)
					})
				})

				r.Routes("/applications", func(r chi.Router) {
					r.Get("/", s.execGetApplications)
					r.Get("/{id}", s.execGetApplication)
				})

				// NOTE json payload here should probably be able to delete an announcement
				// TODO remove note above when that's implemented
				r.Post("/announcements", s.execAddAnnouncement)
			})
		*/
	})

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
