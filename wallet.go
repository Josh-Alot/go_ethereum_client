package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
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

func createKeystore() *keystore.KeyStore {
	return keystore.NewKeyStore(accountsDir, keystore.StandardScryptN, keystore.StandardScryptP)
}
