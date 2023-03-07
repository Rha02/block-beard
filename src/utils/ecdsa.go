package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

// Signature is a struct for a signature.
type Signature struct {
	R, S *big.Int
}

// ToString() returns a string representation of the signature.
func (s *Signature) ToString() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

// String2BigIntTuple() converts a string to a tuple of two big.Ints.
func String2BigIntTuple(s string) (*big.Int, *big.Int) {
	x := new(big.Int)
	y := new(big.Int)

	x.SetString(s[:len(s)/2], 16)
	y.SetString(s[len(s)/2:], 16)

	return x, y
}

// PublicKeyFromString() converts a string to a public key.
func PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := String2BigIntTuple(s)
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
}

// PrivateKeyFromString() converts a string to a private key.
func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	d := new(big.Int)
	d.SetString(s, 16)
	return &ecdsa.PrivateKey{PublicKey: *publicKey, D: d}
}

func SignatureFromString(s string) *Signature {
	x, y := String2BigIntTuple(s)
	return &Signature{R: x, S: y}
}
