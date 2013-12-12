package sign

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"math/big"
	"testing"
)

// TODO deal with leading 0s, aka:
//var data []byte=[]byte("\000data")
var data []byte = []byte("data")

func TestSign(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	sig, err := Sign(key, data)
	if err != nil {
		t.Errorf("Failed to sign: %s", err)
	}
	if !CheckSig(&key.PublicKey, data, sig) {
		t.Errorf("Failed to match sig")
	}
}

func TestBlindSign(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	blinded, unblinder := Blind(&key.PublicKey, data)
	sig := BlindSign(key, blinded)
	unblind := Unblind(&key.PublicKey, sig, unblinder)
	if !CheckBlindSig(&key.PublicKey, data, unblind) {
		t.Errorf("Failed to match sig")
	}
}

func TestBlindSignDoubleValidate(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	blinded, _ := Blind(&key.PublicKey, data)
	sig := BlindSign(key, blinded)
	if CheckBlindSig(&key.PublicKey, blinded, sig) {
		t.Errorf("sig matched for blinded data! Testing that you can't cast a blinded ballot, in addition to the unblinded one. This test fails, since it passes validation.")
	}
}

func TestRawSign(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	sig := BlindSign(key, data)
	if !CheckBlindSig(&key.PublicKey, data, sig) {
		t.Errorf("Failed to match sig")
	}
}

func TestRaw(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)

	m := new(big.Int).SetBytes(data)
	encrypted := encrypt(new(big.Int), &key.PublicKey, m)

	sig := BlindSign(key, encrypted.Bytes())
	if !bytes.Equal(data, sig) {
		t.Errorf("Failed to match sig")
	}
}

func TestVeryRaw(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)

	m := new(big.Int).SetBytes(data)
	encrypted := encrypt(new(big.Int), &key.PublicKey, m)

	m2, err := decrypt(rand.Reader, key, encrypted)

	if err != nil {
		t.Errorf("Failed to decrypt: %s", err)
	}

	if m2.Cmp(m) != 0 {
		t.Errorf("Failed to match sig")
	}
}

func TestVeryRaw2(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)

	encrypted := new(big.Int).SetBytes(data)
	m, err := decrypt(rand.Reader, key, encrypted)
	encrypted2 := encrypt(new(big.Int), &key.PublicKey, m)

	if err != nil {
		t.Errorf("Failed to decrypt: %s", err)
	}
	if encrypted2.Cmp(encrypted) != 0 {
		t.Errorf("Failed to match sig")
	}
}

func TestVeryRaw2Bytes(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)

	encrypted := new(big.Int).SetBytes(data)
	m, err := decrypt(rand.Reader, key, encrypted)
	encrypted2 := encrypt(new(big.Int), &key.PublicKey, m)

	if err != nil {
		t.Errorf("Failed to decrypt: %s", err)
	}
	if !bytes.Equal(encrypted.Bytes(), encrypted2.Bytes()) {
		t.Errorf("Failed to match sig")
	}
}
