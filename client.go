package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"crypto/rsa"
	"fmt"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/msg/msgs"
	"github.com/Craig-Macomber/election/sign"
	"net"
)

var ballotKey *rsa.PublicKey

func init() {
	ballotKey = keys.UnpackKey(keys.LoadKey("publicData/ballotKey"))
}

const maxLength = 4096

func main() {
    msgKey := keys.LoadPrivateKey("clientPrivateData/voterKey")
    voterKey := keys.UnpackPrivateKey(msgKey)
	
	keySig,err:=GetKeySig([]byte(keys.StringKey(msgKey.PublicKey)))
	if err != nil {
		panic(err)
	}
	
	// Construct ballot
	// TODO: prompt user or read from file
	ballot := []byte("ballot!!")
	
	ballotSig,err:=GetBallotSig(voterKey,keySig,ballot)
	if err != nil {
		panic(err)
	}
	
	vote,err:=SubmitBallot(ballot,ballotSig)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Cast ballot '%s' as %s",ballot,vote)
}

// Get signature from Voter List Server to prove our public key is valid
func GetKeySig(key []byte) ([]byte,error){
    // connect
    conn, err := net.Dial("tcp", "localhost"+msg.Service)
	if err != nil {
		return nil,err
	}
	// send request
	err = msg.WriteBlock(conn, msg.KeySignatureRequest, key)
	if err != nil {
		return nil,err
	}
	// Read response
	t, err := msg.ReadType(conn)
	if err != nil {
		return nil,err
	}
	if t != msg.KeySignatureResponse {
		return nil,fmt.Errorf("invalid response type %d",t)
	}
	sig, err := msg.ReadBlock(conn, maxLength)
	conn.Close()
	return sig,err
}

// Get signature for the ballot from ballot server
func GetBallotSig(voterKey *rsa.PrivateKey, keySig, ballot []byte) ([]byte,error){
	

	conn, err := net.Dial("tcp", "localhost"+msg.Service)
	if err != nil {
		panic(err)
	}

	var r msgs.SignatureRequest

	// TODO real values
	r.VoterPublicKey = keys.PackKey(&voterKey.PublicKey)
	

	blindedBallot, unblinder := sign.Blind(ballotKey, ballot)
	r.BlindedBallot = blindedBallot
	r.KeySignature=keySig
	voterSignature, err := sign.Sign(voterKey, blindedBallot)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	r.VoterSignature = voterSignature
	data, err := proto.Marshal(&r)
	if err != nil {
		panic(err)
	}
	err = msg.WriteBlock(conn, msg.SignatureRequest, data)
	if err != nil {
		panic(err)
	}
	t, err := msg.ReadType(conn)
	if err != nil {
		panic(err)
	}
	if t != msg.SignatureResponse {
		panic("invalid response type")
	}
	data, err = msg.ReadBlock(conn, maxLength)
	conn.Close()
	if err != nil {
		panic(err)
	}
	var response msgs.SignatureResponse
	err = proto.Unmarshal(data, &response)
	if err != nil {
		fmt.Println("error reading SignatureResponse:", err)
		panic(err)
	}

	if !KeysEqual(r.VoterPublicKey, response.Request.VoterPublicKey) {
		fmt.Println("illegal response from server. Voter keys must match:")
		fmt.Println(r.VoterPublicKey)
		fmt.Println(response.Request.VoterPublicKey)
		panic(err)
	}

	err = msg.ValidateSignatureRequest(response.Request)
	if err != nil {
		panic(err)
	}

	newBlindedBallot := response.Request.BlindedBallot

	if !bytes.Equal(r.BlindedBallot, newBlindedBallot) {
		fmt.Println("Server has proof you voted before!")
		// TODO should store ballots (and blinding factors?) from all attempts
		// so can try and cast old ballot
		panic(err)
	}

	if !sign.CheckBlindSig(ballotKey, newBlindedBallot, response.BlindedBallotSignature) {
		fmt.Println("illegal response from server. Signature in response is invalid:", ballotKey, newBlindedBallot, response.BlindedBallotSignature)
		panic(err)
	}

	sig := sign.Unblind(ballotKey, response.BlindedBallotSignature, unblinder)

    return sig,nil
}

func KeysEqual(a, b *msgs.PublicKey) bool {
	if *a.E != *b.E {
		return false
	}
	return bytes.Equal(a.N, b.N)
}

func SubmitBallot(ballot, sig []byte) (*msgs.VoteResponse,error ){
	voteKey := keys.UnpackKey(keys.LoadKey("publicData/voteKey"))

	fmt.Printf("Casting Ballot: %s\n", ballot)

	conn, err := net.Dial("tcp", "localhost"+msg.Service)
	if err != nil {
		return nil,err
	}

	var vote msgs.Vote
	vote.Ballot = ballot
	vote.BallotSignature = sig

	// redundant sanity check signature
	err = msg.ValidateVote(ballotKey, &vote)
	if err != nil {
		return nil,err
	}
	data, err := proto.Marshal(&vote)
	if err != nil {
		return nil,err
	}
	err = msg.WriteBlock(conn, msg.Vote, data)
	if err != nil {
		return nil,err
	}
	t, err := msg.ReadType(conn)
	if err != nil {
		return nil,err
	}
	if t != msg.VoteResponse {
		return nil,fmt.Errorf("invalid response type")
	}
	data, err = msg.ReadBlock(conn, maxLength)
	conn.Close()
	if err != nil {
		return nil,err
	}

	var response msgs.VoteResponse
	err = proto.Unmarshal(data, &response)
	if err != nil {
		fmt.Println("error reading VoteResponse:", err)
		return nil,err
	}

	b := response.BallotEntry
	s := response.BallotEntrySignature

	if !sign.CheckSig(voteKey, b, s) {
		err = fmt.Errorf("illegal vote response from server. Signature in request is invalid")
		return nil,err
	}

	fmt.Printf("Got signed BallotEntry for: %s\n", ballot)
	return &response,nil
}
