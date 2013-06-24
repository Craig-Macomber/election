package keys

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rsa"
	"github.com/Craig-Macomber/election/msg/msgs"
	"io/ioutil"
	"math/big"
	"os"
)

func PackKey(k *rsa.PublicKey) *msgs.PublicKey {
	var key msgs.PublicKey
	key.N = k.N.Bytes()
	tmp := int64(k.E)
	key.E = &tmp
	return &key
}

func PackPrivateKey(k *rsa.PrivateKey) *msgs.PrivateKey {
	var key msgs.PrivateKey
	key.PublicKey = PackKey(&k.PublicKey)
	key.D = k.D.Bytes()
	for _, p := range k.Primes {
		key.Primes = append(key.Primes, p.Bytes())
	}
	return &key
}

func UnpackKey(k *msgs.PublicKey) *rsa.PublicKey {
	var key rsa.PublicKey
	key.N = new(big.Int)
	key.N.SetBytes(k.N)
	key.E = int(*k.E)
	return &key
}

func UnpackPrivateKey(k *msgs.PrivateKey) *rsa.PrivateKey {
	var key rsa.PrivateKey
	key.PublicKey = *UnpackKey(k.PublicKey)
	key.D = new(big.Int)
	key.D.SetBytes(k.D)
	for _, p := range k.Primes {
		key.Primes = append(key.Primes, new(big.Int).SetBytes(p))
	}
	return &key
}

func LoadKey(path string) *msgs.PublicKey {
	// open input file
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	data, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	var k msgs.PublicKey
	err = proto.Unmarshal(data, &k)
	if err != nil {
		panic(err)
	}
	return &k
}

func LoadPrivateKey(path string) *msgs.PrivateKey {
	// open input file
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	data, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	var k msgs.PrivateKey
	err = proto.Unmarshal(data, &k)
	if err != nil {
		panic(err)
	}
	return &k
}
