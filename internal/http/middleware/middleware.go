package middleware

type contextKey string

func (c contextKey) String() string {
	return "middleware context key " + string(c)
}
