package status

type Status int

const (
	Offseaon = iota + 1
	StatusApplicationsOpen
	StatusApplicationsClosed
	StatusLive
	StatusPostseason
)
