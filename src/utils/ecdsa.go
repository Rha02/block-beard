package utils

import (
	"fmt"
	"math/big"
)

// Signature is a struct for a signature.
type Signature struct {
	R, S *big.Int
}

// ToString() returns a string representation of the signature.
func (s *Signature) ToString() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
