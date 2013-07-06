// Package ballot provides all needed handlers to run a "Ballot Server"
// A Ballot Server signes (blinded) ballot. It will sign at most one ballot per voter public key
// A request for additional singings will be replied to with a copy of the response to the first request
// Since the responses contain the request, which is signed by the voter, this serves as proof they
// previously submitted a request.
package ballot

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
	"sync"
)

var ballotPrivateKey *rsa.PrivateKey
var voterListKey *rsa.PublicKey

func init() {
	ballotPrivateKey = keys.UnpackPrivateKey(config.LoadServerKey("ballotKey"))
	config := config.Load()
	voterListKey = keys.UnpackKey(config.VoterListServer.Key)
}

// Map of voter public key -> msg.SignatureResponse containing the signed ballot
// This serves as proof that someone with the corresponding private key tried to vote
// This is used to prove at least len(signatureResponses) voters (This must be >= the number of final ballots!)
// As well as to prevent multiple voting (a request to vote with a different ballot returns
// the saved previous response as proof they already had a ballot signed. (If they failed to receive it previously, they well get it again)
// Note: when counting, only valid (belonging to real voters) public keys should be counted
var signatureResponses map[string][]byte = map[string][]byte{}
var responseLock sync.Mutex

func getResponse(keyString string, r *msgs.SignatureRequest) []byte {
	responseLock.Lock()
	responseData, ok := signatureResponses[keyString]
	if !ok {
		var response msgs.SignatureResponse
		response.Request = r
		response.BlindedBallotSignature = sign.BlindSign(ballotPrivateKey, r.BlindedBallot)
		responseData, _ = proto.Marshal(&response)
		signatureResponses[keyString] = responseData
		fmt.Printf("Signed ballot for %s\n", *keys.UnpackKey(r.VoterPublicKey))
	}
	responseLock.Unlock()
	return responseData
}

func HandelSignatureRequest(data []byte, c net.Conn) {
	var r msgs.SignatureRequest
	err := proto.Unmarshal(data, &r)
	if err != nil {
		fmt.Println("server error reading SignatureRequest:", err)
		server.ConnectionError(c)
		return
	}

	err = msg.ValidateSignatureRequest(&r)
	if err != nil {
		fmt.Println(err)
		return
	}

	keyString := keys.StringKey(r.VoterPublicKey)

	if !sign.CheckSig(voterListKey, []byte(keyString), r.KeySignature) {
		fmt.Println("SignatureRequest's KeySignature Signature is invalid")
		server.ConnectionError(c)
		return
	}

	responseData := getResponse(keyString, &r)

	server.SendBlock(msg.SignatureResponse, responseData, c)
}
