// Some helpers for converting keys between formats
package keys

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/rsa"
	"github.com/Craig-Macomber/election/msg/msgs"
	"io/ioutil"
	"math/big"
	"os"
)

// Serialize a public key to a string (useful for using keys as map keys)
func StringKey(k *msgs.PublicKey) string {
	keyBytes, err := proto.Marshal(k)
	if err != nil {
		panic(err)
	}
	return string(keyBytes)
}

// rsa -> msg
func PackKey(k *rsa.PublicKey) *msgs.PublicKey {
	var key msgs.PublicKey
	key.N = k.N.Bytes()
	tmp := int64(k.E)
	key.E = &tmp
	return &key
}

// rsa -> msg
func PackPrivateKey(k *rsa.PrivateKey) *msgs.PrivateKey {
	var key msgs.PrivateKey
	key.PublicKey = PackKey(&k.PublicKey)
	key.D = k.D.Bytes()
	for _, p := range k.Primes {
		key.Primes = append(key.Primes, p.Bytes())
	}
	return &key
}

// msg -> rsa
func UnpackKey(k *msgs.PublicKey) *rsa.PublicKey {
	var key rsa.PublicKey
	key.N = new(big.Int)
	key.N.SetBytes(k.N)
	key.E = int(*k.E)
	return &key
}

// msg -> rsa
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

// trivial read all helper
func LoadBytes(path string) []byte{
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
	return data
}

// Load key from file
func LoadKey(path string) *msgs.PublicKey {
	data:=LoadBytes(path)
	var k msgs.PublicKey
	err := proto.Unmarshal(data, &k)
	if err != nil {
		panic(err)
	}
	return &k
}

// Load key from file
func LoadPrivateKey(path string) *msgs.PrivateKey {
	data:=LoadBytes(path)
	var k msgs.PrivateKey
	err := proto.Unmarshal(data, &k)
	if err != nil {
		panic(err)
	}
	return &k
}
