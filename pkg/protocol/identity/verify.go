package identity

import (
	"github.com/itspablomontes/fleming/pkg/protocol/crypto"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type Verifier interface {
	VerifySignature(message, signatureHex string, address types.WalletAddress) bool
}

type DefaultVerifier struct{}

func NewVerifier() *DefaultVerifier {
	return &DefaultVerifier{}
}

func (v *DefaultVerifier) VerifySignature(message, signatureHex string, address types.WalletAddress) bool {
	return crypto.VerifySignature(message, signatureHex, address.String())
}

func VerifySIWE(opts SIWEOptions, signatureHex string) (bool, error) {
	if err := opts.Validate(); err != nil {
		return false, err
	}

	message := BuildSIWEMessage(opts)
	verifier := NewVerifier()
	return verifier.VerifySignature(message, signatureHex, opts.Address), nil
}
