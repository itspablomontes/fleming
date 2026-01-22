package crypto

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(message string, signatureHex string, addressHex string) bool {
	sig, err := hexutil.Decode(signatureHex)
	if err != nil {
		return false
	}

	if len(sig) != 65 {
		return false
	}

	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	hash := crypto.Keccak256([]byte(prefix))

	pubKeyBytes, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return false
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	return strings.EqualFold(recoveredAddr.Hex(), addressHex)
}
