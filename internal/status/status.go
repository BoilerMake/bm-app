package status

type Status int

const (
	Offseaon = iota + 1
	ApplicationsOpen
	ApplicationsClosed
	Live
	Postseason
)
