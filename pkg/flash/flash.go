package flash

const (
	Error = iota
	Warning
	Info
	Success
)

// A Flash holds a message and its type
type Flash struct {
	Type    int
	Message string
}
