package msg

import (
	"crypto/rsa"
	"fmt"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg/msgs"
	"github.com/Craig-Macomber/election/sign"
)

// Does not check keySignature
func ValidateSignatureRequest(m *msgs.SignatureRequest) error {
	newBlindedBallot := m.BlindedBallot
	newSig := m.VoterSignature
	key := keys.UnpackKey(m.VoterPublicKey)
	if !sign.CheckSig(key, newBlindedBallot, newSig) {
		return fmt.Errorf("SignatureRequest's VoterSignature Signature is invalid")
	}
	return nil
}

func ValidateVote(ballotKey *rsa.PublicKey, m *msgs.Vote) error {
	ballot := m.Ballot
	sig := m.BallotSignature
	if !sign.CheckBlindSig(ballotKey, ballot, sig) {
		return fmt.Errorf("Vote's BallotSignature Signature is invalid")
	}
	return nil
}
