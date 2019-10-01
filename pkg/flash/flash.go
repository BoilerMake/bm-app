package flash

const (
	FlashError = iota
	FlashWarning
	FlashInfo
	FlashSuccess
)

// A Flash holds a message and its type
type Flash struct {
	Type    int
	Message string
}
