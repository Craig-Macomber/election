package msg

type Type uint16

const (
	SignatureRequest Type = iota
	SignatureResponse
	Vote
	VoteResponse
)
