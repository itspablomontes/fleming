package identity

import (
	"fmt"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type SIWEOptions struct {
	Address types.WalletAddress
	Domain  string

	URI string

	Nonce string

	ChainID int

	IssuedAt time.Time

	ExpirationTime *time.Time
	Statement      string
}

const DefaultStatement = "Sign in to Fleming for secure access to your medical data."

func BuildSIWEMessage(opts SIWEOptions) string {
	issuedAt := opts.IssuedAt
	if issuedAt.IsZero() {
		issuedAt = time.Now().UTC()
	}

	statement := opts.Statement
	if statement == "" {
		statement = DefaultStatement
	}

	msg := fmt.Sprintf(`%s wants you to sign in with your Ethereum account:
%s

%s

URI: %s
Version: 1
Chain ID: %d
Nonce: %s
Issued At: %s`,
		opts.Domain,
		opts.Address.String(),
		statement,
		opts.URI,
		opts.ChainID,
		opts.Nonce,
		issuedAt.Format(time.RFC3339),
	)

	if opts.ExpirationTime != nil {
		msg += fmt.Sprintf("\nExpiration Time: %s", opts.ExpirationTime.Format(time.RFC3339))
	}

	return msg
}

func (opts SIWEOptions) Validate() error {
	var errs types.ValidationErrors

	if opts.Address.IsEmpty() {
		errs.Add("address", "wallet address is required")
	}

	if opts.Domain == "" {
		errs.Add("domain", "domain is required")
	}

	if opts.URI == "" {
		errs.Add("uri", "URI is required")
	}

	if opts.Nonce == "" {
		errs.Add("nonce", "nonce is required")
	}

	if opts.ChainID <= 0 {
		errs.Add("chainId", "chain ID must be positive")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
