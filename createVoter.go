package main

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rand"
	"crypto/rsa"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg/msgs"
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

func makeVoter(publicPath, privatePath, name string) {
	var voter msgs.VoterData
	voter.Name = &name

	ballotKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	k := keys.PackPrivateKey(ballotKey)
	voter.Key = k
	private, err := proto.Marshal(k.PublicKey)
	if err != nil {
		panic(err)
	}
	doFile(publicPath, private)

	// Write out file
	voterData, err := proto.Marshal(&voter)
	if err != nil {
		panic(err)
	}
	doFile(privatePath, voterData)
}

func main() {
	electionPath := "demoElection/"
	voterPublicKeyDir := electionPath + "voterPublicKeys/"
	privateOutput := electionPath + "voterPrivate/"

	names := []string{"TestVoter1", "TestVoter2", "TestVoter3"}
	for _, name := range names {
		makeVoter(voterPublicKeyDir+name, privateOutput+name, name)
	}
}
