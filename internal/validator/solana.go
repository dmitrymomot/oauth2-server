package validator

import (
	"fmt"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/common"
)

// ValidateSolanaWalletAddr validates a Solana wallet address.
// Returns an error if the address is invalid, nil otherwise.
func ValidateSolanaWalletAddr(addr string) error {
	if addr == "" {
		return fmt.Errorf("wallet address is empty")
	}

	d, err := base58.Decode(addr)
	if err != nil {
		return fmt.Errorf("invalid wallet address: %w", err)
	}

	if len(d) != common.PublicKeyLength {
		return fmt.Errorf("invalid wallet address length: %d", len(d))
	}

	if _, err := new(edwards25519.Point).SetBytes(d); err != nil {
		return fmt.Errorf("invalid wallet address: %w", err)
	}

	return nil
}
