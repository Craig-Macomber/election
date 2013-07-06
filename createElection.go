package main

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rand"
	"crypto/rsa"
	"github.com/Craig-Macomber/election/keys"
	"github.com/Craig-Macomber/election/msg"
	"github.com/Craig-Macomber/election/msg/msgs"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
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

func makeServer(publicDst **msgs.Server, privatePath string) {
	wg.Add(1)
	go func() {
		var server msgs.Server
		*publicDst = &server
		address := ("localhost" + msg.Service)
		server.Address = &address
		ballotKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		k := keys.PackPrivateKey(ballotKey)
		server.Key = k.PublicKey
		private, err := proto.Marshal(k)
		if err != nil {
			panic(err)
		}
		doFile(privatePath, private)
		wg.Done()
	}()
}

var wg sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(4)
	electionPath := "demoElection/"
	voterPublicKeyDir := electionPath + "voterPublicKeys/"
	privateOutput := electionPath + "serverPrivate/"
	configDst := electionPath + "config"
	var config msgs.ElectionConfig

	// Generate server keys
	makeServer(&config.BallotServer, privateOutput+"ballotKey")
	makeServer(&config.VoteServer, privateOutput+"voteKey")
	makeServer(&config.VoterListServer, privateOutput+"voterListKey")
	makeServer(&config.FinalVoteSetServer, privateOutput+"finalVoteSetKey")
	makeServer(&config.FinalSignatureRequestSetServer, privateOutput+"finalSignatureRequestSetKey")

	// Fill in voters
	infos, err := ioutil.ReadDir(voterPublicKeyDir)
	if err != nil {
		panic(err)
	}
	for _, info := range infos {
		name := info.Name()
		var v msgs.Voter
		v.Key = keys.LoadKey(voterPublicKeyDir + name)
		v.Name = &name
		config.Voters = append(config.Voters, &v)
	}

	// Fill in description
	s := "A test election"
	config.BallotDescription = &s

	// Wait for server keys
	wg.Wait()

	// Write out file
	configData, err := proto.Marshal(&config)
	if err != nil {
		panic(err)
	}
	doFile(configDst, configData)
}
