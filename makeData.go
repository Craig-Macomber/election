package main

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rand"
	"crypto/rsa"
	"github.com/Craig-Macomber/election/keys"
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

func makeKey(publicPath, privatePath string) {
	wg.Add(1)
	go func() {
		ballotKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		k := keys.PackPrivateKey(ballotKey)
		private, err := proto.Marshal(k)
		if err != nil {
			panic(err)
		}
		public, err := proto.Marshal(k.PublicKey)
		if err != nil {
			panic(err)
		}
		doFile(privatePath, private)
		doFile(publicPath, public)
		wg.Done()
	}()
}

var wg sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(4)
	makeKey("publicData/ballotKey", "serverPrivateData/ballotKey")
	makeKey("publicData/voterKey", "clientPrivateData/voterKey")
	makeKey("publicData/voteKey", "serverPrivateData/voteKey")
	makeKey("publicData/voterListKey", "serverPrivateData/voterListKey")
	wg.Wait()
}
