// Package vote provides all needed handlers to run a "Vote Server"
// A Vote server does the actual collection of all the votes (signed ballots).
// Each vote must be signed by the ballot server.
package vote

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rsa"
	"fmt"
	"github.com/Craig-Macomber/election/config"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/msg/msgs"
	"github.com/Craig-Macomber/election/server"
	"github.com/Craig-Macomber/election/sign"
	"net"
)

var votePrivateKey *rsa.PrivateKey
var ballotKey *rsa.PublicKey

func init() {
	votePrivateKey = keys.UnpackPrivateKey(config.LoadServerKey("voteKey"))
	config := config.Load()
	ballotKey = keys.UnpackKey(config.BallotServer.Key)
}

// All the votes
// The index is the id
var votes []msgs.Vote = make([]msgs.Vote, 0)

// A set of all ballots in votes, maps to id
var ballotSet map[string]int = map[string]int{}

func stringKey(v *msgs.Vote) string {
	return string(v.Ballot)
}

func HandelSignatureRequest(data []byte, c net.Conn) {
	var v msgs.Vote
	err := proto.Unmarshal(data, &v)
	if err != nil {
		fmt.Println("server error reading Vote:", err)
		server.ConnectionError(c)
		return
	}

	err = msg.ValidateVote(ballotKey, &v)
	if err != nil {
		fmt.Println(err)
		return
	}

	keyString := stringKey(&v)
	// TODO: needs a lock
	id, ok := ballotSet[keyString]
	if !ok {
		id = len(votes)
		votes = append(votes, v)
		ballotSet[keyString] = id
		fmt.Printf("Got Vote %d: %s\n", id, v.Ballot)
	}

	var response msgs.VoteResponse
	var ballotEntry msgs.BallotEntry
	tmp := uint64(id)
	ballotEntry.Id = &tmp
	ballotEntry.Ballot = votes[id].Ballot

	ballotBytes, err := proto.Marshal(&ballotEntry)

	response.BallotEntry = ballotBytes
	response.BallotEntrySignature, err = sign.Sign(votePrivateKey, ballotBytes)

	responseData, err := proto.Marshal(&response)
	server.SendBlock(msg.VoteResponse, responseData, c)
}
