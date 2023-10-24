package wallet

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/common"
)

func (w *Wallet) ListAccounts() ([]string, error) {
	var err error
	var result []string
	c, err := rpc.Dial(w.Host)
	if err != nil {
		return nil, fmt.Errorf("rpc dial error: %v", err)
	}
	err = c.Call(&result, "eth_accounts")
	if err != nil {
		return nil, fmt.Errorf("eth.accounts error: %v", err)
	}
	return result, nil
}

func (w *Wallet) GetBalance(account string) (float64, error) {
	balance, err := w.client.BalanceAt(context.Background(), common.HexToAddress(account), nil)
	if err != nil {
		return 0, err
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	result, _ := ethValue.Float64()
	return result, nil
}

func (w *Wallet) CreatNewAccount(password string) (string, error) {
	var err error
	var result string
	c, err := rpc.Dial(w.Host)
	if err != nil {
		return "", fmt.Errorf("rpc dial error: %v", err)
	}
	err = c.Call(&result, "personal_newAccount", password)
	if err != nil {
		return "", err
	}
	return result, nil
}
