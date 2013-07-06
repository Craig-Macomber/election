// Package voterlist provides all needed handlers to run a "Voter List Server"
// A "Voter List Server" provides a list of registered voters' public keys.
// It will sign valid voter public keys with its voterListKey.
//
// This server is completely static and stateless. It only serves static content, which is all public.
package voterlist

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/Craig-Macomber/election/config"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/server"
	"github.com/Craig-Macomber/election/sign"
	"net"
)

func init() {

	// Load the singing key
	// TODO: load from election config file
	voterListKey := keys.UnpackPrivateKey(config.LoadServerKey("voterListKey"))
	config := config.Load()
	for _, voter := range config.Voters {
		voterKey, err := proto.Marshal(voter.Key)
		if err != nil {
			panic(err)
		}
		sig, err := sign.Sign(voterListKey, voterKey)
		if err != nil {
			panic(err)
		}
		voterKeys[string(voterKey)] = sig
		fmt.Printf("Added voter %s.\n", *voter.Name)
	}
}

// A map of voter public keys -> signatures for them (signed with voterListKey)
var voterKeys map[string][]byte = map[string][]byte{}

func HandelSignatureRequest(data []byte, c net.Conn) {
	keyString := string(data)
	responseData, ok := voterKeys[keyString]
	if !ok {
		fmt.Println("SignatureRequest with unknown key:", keyString)
		server.ConnectionError(c)
		return
	}
	server.SendBlock(msg.KeySignatureResponse, responseData, c)
}

// TODO: add way go get full list of voters. The signature of this list (signed with voterListKey)
// should be included in the election description file
// (or for small elections, the full list could be included)
