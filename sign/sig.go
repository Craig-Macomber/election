package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	"math/big"
)

var hashType crypto.Hash = crypto.SHA256

func Hash(data []byte) []byte {
	h := hashType.New()
	h.Write(data)
	return h.Sum(make([]byte, 0))
}

func Sign(key *rsa.PrivateKey, data []byte) (sig []byte, err error) {
	hashResult := Hash(data)
	sig, err = rsa.SignPKCS1v15(rand.Reader, key, hashType, hashResult)
	return
}

func CheckSig(key *rsa.PublicKey, data, sig []byte) bool {
	hashResult := Hash(data)
	err := rsa.VerifyPKCS1v15(key, hashType, hashResult, sig)
	return err == nil
}

func BlindSign(key *rsa.PrivateKey, data []byte) []byte {
	c := new(big.Int).SetBytes(data)
	m, err := decrypt(rand.Reader, key, c)
	if err != nil {
		// TODO: handel errors that be caused by bad user input
		panic(err)
	}
	return m.Bytes()
}

func Blind(key *rsa.PublicKey, data []byte) (blindedData, unblinder []byte) {
	blinded, unblinderBig, err := blind(rand.Reader, key, new(big.Int).SetBytes(data))
	if err != nil {
		panic(err)
	}
	return blinded.Bytes(), unblinderBig.Bytes()
}

func Unblind(key *rsa.PublicKey, blindedSig, unblinder []byte) []byte {
	m := new(big.Int).SetBytes(blindedSig)
	unblinderBig := new(big.Int).SetBytes(unblinder)
	m.Mul(m, unblinderBig)
	m.Mod(m, key.N)
	return m.Bytes()
}

func CheckBlindSig(key *rsa.PublicKey, data, sig []byte) bool {
	m := new(big.Int).SetBytes(data)
	bigSig := new(big.Int).SetBytes(sig)
	c := encrypt(new(big.Int), key, bigSig)
	return m.Cmp(c) == 0
}