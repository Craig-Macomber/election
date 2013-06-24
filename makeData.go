package main

import (
	"github.com/Craig-Macomber/election/keys"

	"code.google.com/p/goprotobuf/proto"

	"crypto/rand"
	"crypto/rsa"
	"os"
)

func doFile(path string, data []byte) {
	// open output file
	fo, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = fo.Write(data)
	if err != nil {
		panic(err)
	}
}

func main() {
	ballotKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	k := keys.PackPrivateKey(ballotKey)
	private, err := proto.Marshal(k)
	if err != nil {
		panic(err)
	}
	public, err := proto.Marshal(k.PublicKey)
	if err != nil {
		panic(err)
	}
	doFile("serverPrivateData/ballotKey", private)
	doFile("publicData/ballotKey", public)

	voterKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	k = keys.PackPrivateKey(voterKey)
	private, err = proto.Marshal(k)
	if err != nil {
		panic(err)
	}
	public, err = proto.Marshal(k.PublicKey)
	if err != nil {
		panic(err)
	}
	doFile("clientPrivateData/voterKey", private)
	doFile("publicData/voterKey", public)

	voteKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	k = keys.PackPrivateKey(voteKey)
	private, err = proto.Marshal(k)
	if err != nil {
		panic(err)
	}
	public, err = proto.Marshal(k.PublicKey)
	if err != nil {
		panic(err)
	}
	doFile("serverPrivateData/voteKey", private)
	doFile("publicData/voteKey", public)
}
