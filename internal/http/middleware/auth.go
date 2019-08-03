package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

// WithJWT provides a middleware that attempts to decode a JWT from a request's
// cookies.  It sets two fields in the requests context, one with the actual
// JWT and the other with an error that may have come from decoding that JWT.
// When checking if a request has a valid JWT you should always check that the
// error field is nil, and that the token field is not nil (in that order). For
// an example of that check out getClaimsFromCtx in internal/http/api/handler.go

// func WithJWT(next http.Handler) http.Handler {
// 	sessCookie, ok := os.LookupEnv("SESS_COOKIE_NAME")
// 	if !ok {
// 		log.Fatalf("environment variable not set: %v", "JWT_COOKIE_NAME")
// 	}
//
// 	JWTSigningKeyString, ok := os.LookupEnv("JWT_SIGNING_KEY")
// 	if !ok {
// 		log.Fatalf("environment variable not set: %v", "JWT_SIGNING_KEY")
// 	}
// 	JWTSigningKey := []byte(JWTSigningKeyString)
//
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		token, err := getToken(r, jwtCookie, JWTSigningKey)
//
// 		ctx := r.Context()
// 		ctx = context.WithValue(ctx, JWTCtxKey, token)
// 		ctx = context.WithValue(ctx, JWTErrorCtxKey, err)
//
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
//
// 	return http.HandlerFunc(fn)
// }

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("SESS_COOKIE_NAME")
	store = sessions.NewCookieStore(key)
)

func sessionMiddleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			// log.WithError(err).Error("bad session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "session", session))
		h(w, r)
	}
}

// // DecodeAndModify provides a way to
// func Decode(w http.ResponseWriter, r *http.Request) http.Handler {
//     //get reference to cookie if set
//     if cookie, err := r.Cookie("session-cookie-name"); err == nil {
//
//         value := make(map[string]string)
//         //use Decode to get the value from the cookie
//         if err = s.Decode("session-cookie-name", cookie.Value, &value); err == nil {
//             //modify the value in some way
//
//         }
//     }
// }
// getToken decodes a JWT.  This func is separated not only to shorten the
// WithJWT func but more importantly to make the whole "return token, err"
// process more clear.
// func getToken(r *http.Request, jwtCookie string, JWTSigningKey []byte) (token *jwt.Token, err error) {
// 	cookie, err := r.Cookie(jwtCookie)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	token, err = jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
// 		// Don't forget to validate the alg is what you expect:
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
//
// 		return JWTSigningKey, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		// TODO refresh tokens here and put into request cookies here, keep in mind
// 		// they *may* be rewritten later on in the request chain, but that shouldn't
// 		// really matter as long as the resulting token is still valid.
// 		return token, err
// 	}
// }
