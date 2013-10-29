package main

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"crypto/rsa"
	"crypto/rand"
	"fmt"
	"github.com/Craig-Macomber/election/config"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/msg/msgs"
	"github.com/Craig-Macomber/election/sign"
	"net"
	"encoding/base64"
)

var ballotKey *rsa.PublicKey
var voteKey *rsa.PublicKey

const maxLength = 4096

func LoadVoterData(path string) *msgs.VoterData {
	data := keys.LoadBytes(path)
	var voterData msgs.VoterData
	err := proto.Unmarshal(data, &voterData)
	if err != nil {
		panic(err)
	}
	return &voterData
}


func load(name string) (*msgs.ElectionConfig, *msgs.VoterData) {
    privateInfo := LoadVoterData("demoElection/voterPrivate/" + name)
	configBytes := config.LoadBytes()
	privateInfo.ElectionConfig = configBytes
	config := config.Unpack(privateInfo.ElectionConfig)
	return config,privateInfo
}

func ConfigHash(privateInfo *msgs.VoterData) []byte {
    return sign.Hash(privateInfo.ElectionConfig)
}

type SignatureStatus uint8

const (
	Missing SignatureStatus = iota
	Invalid
	Valid
)

func CheckKeySig(privateInfo *msgs.VoterData) SignatureStatus{
    if privateInfo.KeySignature==nil {
        return Missing
    }
    config := config.Unpack(privateInfo.ElectionConfig)
    voterListKey := keys.UnpackKey(config.VoterListServer.Key)
    publicKey:=PublicKey(privateInfo)
    if sign.CheckSig(voterListKey,publicKey,privateInfo.KeySignature) {
        return Valid
    }
    return Invalid
}

func fillInfo(privateInfo *msgs.VoterData) {
    //config := config.Unpack(privateInfo.ElectionConfig)
    for CheckKeySig(privateInfo)!=Valid {
        // TODO: prompt user here? (Maybe have an interactive mode, only try N times if not interactive?
        fmt.Println("Attempting to get valid keySignature from voter list server...")
        var err error
        privateInfo.KeySignature, err = GetKeySig(PublicKey(privateInfo))
	    if err != nil {
		    fmt.Println(err)
	    } else {
	        fmt.Println("Got keySignature from voter list server. Checking...")
	    }
    }
    fmt.Println("Got valid keySignature from voter list server.")
}

func PublicKey(privateInfo *msgs.VoterData) []byte {
    return []byte(keys.StringKey(privateInfo.Key.PublicKey))
}

// Prefix ballot with random 64 bits to make ballot unique.
// Ballots get de-duplicated by the voteServer, so its important
// that clients make their ballot unique.
func PrefixBallot(payload []byte) []byte{
    const randomLength=8
    b:=make([]byte,len(payload)+randomLength)
    _,err:=rand.Read(b[:randomLength])
    if err!=nil {
        panic(err)
    }
    copy(b[randomLength:],payload)
    return b
}

func main() {
	name := "TestVoter1"
	config,privateInfo:=load(name)
	configHash := ConfigHash(privateInfo)
	fmt.Printf("Loaded election config with hash: %s\n",base64.StdEncoding.EncodeToString(configHash))
	
	ballotKey = keys.UnpackKey(config.BallotServer.Key)
	voteKey = keys.UnpackKey(config.VoteServer.Key)

	voterKey := keys.UnpackPrivateKey(privateInfo.Key)
    
    fillInfo(privateInfo)
    //var err error
	//privateInfo.KeySignature, err = GetKeySig(PublicKey(privateInfo))
	//if err != nil {
	//	panic(err)
	//}

	// Construct ballot
	// TODO: prompt user or read from file
	ballot := PrefixBallot([]byte("ballot!!"))

	ballotSig, err := GetBallotSig(voterKey, privateInfo.KeySignature, ballot)
	if err != nil {
		panic(err)
	}

	vote, err := SubmitBallot(ballot, ballotSig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Cast ballot '%s' as %s\n", ballot, vote)
}

// Get signature from Voter List Server to prove our public key is valid
func GetKeySig(key []byte) ([]byte, error) {
	// connect
	conn, err := net.Dial("tcp", "localhost"+msg.Service)
	if err != nil {
		return nil, err
	}
	// send request
	err = msg.WriteBlock(conn, msg.KeySignatureRequest, key)
	if err != nil {
		return nil, err
	}
	// Read response
	t, err := msg.ReadType(conn)
	if err != nil {
		return nil, err
	}
	if t != msg.KeySignatureResponse {
		return nil, fmt.Errorf("KeySignatureResponse: invalid response type %d. Expected: %d", t, msg.KeySignatureResponse)
	}
	sig, err := msg.ReadBlock(conn, maxLength)
	conn.Close()

	// TODO check sig is valid
	return sig, err
}

// Get signature for the ballot from ballot server
func GetBallotSig(voterKey *rsa.PrivateKey, keySig, ballot []byte) ([]byte, error) {

	conn, err := net.Dial("tcp", "localhost"+msg.Service)
	if err != nil {
		panic(err)
	}

	var r msgs.SignatureRequest

	// TODO real values
	r.VoterPublicKey = keys.PackKey(&voterKey.PublicKey)

	blindedBallot, unblinder := sign.Blind(ballotKey, ballot)
	r.BlindedBallot = blindedBallot
	r.KeySignature = keySig
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
		return nil, fmt.Errorf("SignatureResponse: invalid response type %d. Expected: %d", t, msg.SignatureResponse)
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

	return sig, nil
}

func KeysEqual(a, b *msgs.PublicKey) bool {
	if *a.E != *b.E {
		return false
	}
	return bytes.Equal(a.N, b.N)
}

func SubmitBallot(ballot, sig []byte) (*msgs.VoteResponse, error) {

	fmt.Printf("Casting Ballot: %s\n", ballot)

	conn, err := net.Dial("tcp", "localhost"+msg.Service)
	if err != nil {
		return nil, err
	}

	var vote msgs.Vote
	vote.Ballot = ballot
	vote.BallotSignature = sig

	// redundant sanity check signature
	err = msg.ValidateVote(ballotKey, &vote)
	if err != nil {
		return nil, err
	}
	data, err := proto.Marshal(&vote)
	if err != nil {
		return nil, err
	}
	err = msg.WriteBlock(conn, msg.Vote, data)
	if err != nil {
		return nil, err
	}
	t, err := msg.ReadType(conn)
	if err != nil {
		return nil, err
	}
	if t != msg.VoteResponse {
		return nil, fmt.Errorf("invalid response type")
	}
	data, err = msg.ReadBlock(conn, maxLength)
	conn.Close()
	if err != nil {
		return nil, err
	}

	var response msgs.VoteResponse
	err = proto.Unmarshal(data, &response)
	if err != nil {
		fmt.Println("error reading VoteResponse:", err)
		return nil, err
	}

	b := response.BallotEntry
	s := response.BallotEntrySignature

	if !sign.CheckSig(voteKey, b, s) {
		err = fmt.Errorf("illegal vote response from server. Signature in BallotEntry is invalid")
		return nil, err
	}

	fmt.Printf("Got signed BallotEntry for: %s\n", ballot)
	return &response, nil
}
