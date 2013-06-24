package ballot

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rsa"
	"fmt"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/msg/msgs"
	"github.com/Craig-Macomber/election/server"
	"github.com/Craig-Macomber/election/sign"
	"net"
)

var ballotPrivateKey rsa.PrivateKey

func init() {
	voterKey := keys.LoadKey("publicData/voterKey")
	voterKeys[stringKey(voterKey)] = struct{}{}
	ballotPrivateKey = *keys.UnpackPrivateKey(keys.LoadPrivateKey("serverPrivateData/ballotKey"))
}

// Map of voter public key -> msg.SignatureResponse containing the signed ballot
// This serves as proof that someone with the corresponding private key tried to vote
// This is used to prove at least len(signatureResponses) voters (This must be >= the number of final ballots!)
// As well as to prevent multiple voting (a request to vote with a different ballot returns
// the saved previous response as proof they already had a ballot signed. (If they failed to receive it previously, they well get it again)
// Note: when counting, only valid (belonging to real voters) public keys should be counted
var signatureResponses map[string]msgs.SignatureResponse = map[string]msgs.SignatureResponse{}

// A set of voter public keys
var voterKeys map[string]struct{} = map[string]struct{}{}

func stringKey(k *msgs.PublicKey) string {
	keyBytes, err := proto.Marshal(k)
	if err != nil {
		panic(err)
	}
	return string(keyBytes)
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

	keyString := stringKey(r.VoterPublicKey)
	_, ok := voterKeys[keyString]
	if !ok {
		fmt.Println("SignatureRequest with unknown key:", keyString)
		server.ConnectionError(c)
		return
	}

	// TODO: need lock around this
	response, ok := signatureResponses[keyString]
	if !ok {
		response.Request = &r
		response.BlindedBallotSignature = sign.BlindSign(&ballotPrivateKey, r.BlindedBallot)
		signatureResponses[keyString] = response
		fmt.Printf("Signed ballot for %s\n", *keys.UnpackKey(r.VoterPublicKey))
	}

	responseData, err := proto.Marshal(&response)
	server.SendBlock(msg.SignatureResponse, responseData, c)
}
