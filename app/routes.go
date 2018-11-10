package app

import (
	"github.com/go-chi/chi"
)

func (s *server) routes() {
	s.router.Get("/announcements", s.getAnnouncements)

	s.router.Route("/exec", func(r chi.Router) {
		r.Use(s.adminOnly())

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
