// Package voterlist provides all needed handlers to run a "Voter List Server"
// A "Voter List Server" provides a list of registered voters' public keys.
// It will sign valid voter public keys with its voterListKey.
//
// This server is completely static and stateless. It only serves static content, which is all public.
package voterlist

import (
	"fmt"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/server"
	"github.com/Craig-Macomber/election/sign"
	"net"
)

func init() {
	// Load the singing key
	// TODO: load from election config file
	voterListKey := keys.UnpackPrivateKey(keys.LoadPrivateKey("serverPrivateData/voterListKey"))
	
	// Tmp test code to preload a voter public key
	// TODO: load from elsewhere
	voterKey := keys.LoadBytes("publicData/voterKey")
	sig,err:=sign.Sign(voterListKey, voterKey)
	if err!=nil{
	    panic(err)
	}
	voterKeys[string(voterKey)] = sig
}

// A map of voter public keys -> signatures for them (signed with voterListKey)
var voterKeys map[string][]byte = map[string][]byte{}

func HandelSignatureRequest(data []byte, c net.Conn) {
	keyString:=string(data)
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
