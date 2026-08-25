package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

type AccountSummary struct {
	Address string
}

const accountsDir = "accounts/"

func CreateAccount(password string) (*AccountSummary, error) {
	ks := createKeystore()

	newAccount, err := ks.NewAccount(password)
	if err != nil {
		return nil, fmt.Errorf("account creation error: %w", err)
	}

	account := AccountSummary{Address: newAccount.Address.Hex()}

	return &account, nil
}

func ListAccounts() []AccountSummary {
	ks := createKeystore()
	listAccts := ks.Accounts()

	summary := make([]AccountSummary, 0, len(listAccts))
	for _, account := range listAccts {
		summary = append(summary, AccountSummary{Address: account.Address.Hex()})
	}

	return summary
}

// this function only works because other types of schemes are not in the scope of this project
// if we want to give support to any other scheme, such as trezor or ledger
// this feature must be refactored
func LoadPrivateKey(address string, password []byte) (*ecdsa.PrivateKey, error) {
	ks := createKeystore()

	isAddr := common.IsHexAddress(address)
	if !isAddr {
		return nil, fmt.Errorf("hexAddr %s: invalid address", address)
	}

	wantedAddr := common.HexToAddress(address)
	accounts := ks.Accounts()

	for _, acct := range accounts {
		if wantedAddr == acct.Address {
			data, err := os.ReadFile(acct.URL.Path)
			if err != nil {
				return nil, fmt.Errorf("reading keystore error: %w", err)
			}

			key, err := keystore.DecryptKey(data, string(password))
			if err != nil {
				return nil, fmt.Errorf("decryption key error: %w", err)
			}

			return key.PrivateKey, nil
		}
	}

	return nil, fmt.Errorf("load private key: no account found for address %s", address)
}

func createKeystore() *keystore.KeyStore {
	return keystore.NewKeyStore(accountsDir, keystore.StandardScryptN, keystore.StandardScryptP)
}
