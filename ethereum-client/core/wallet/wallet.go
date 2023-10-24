package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"

	"github.com/tyler-smith/go-bip39"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/ethclient"
)

type KeyPairs struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewKeyPair() *KeyPairs {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil
	}
	private, public, err := NewEthWalletByPath("m/44'/60'/0'/0/0", seed)
	if err != nil {
		return nil
	}
	return &KeyPairs{
		privateKey: private,
		publicKey:  public,
	}
}

func (k *KeyPairs) ChainId() int {
	return 3
}

func (k *KeyPairs) ChainParams() *params.ChainConfig {
	return params.RopstenChainConfig
}

func (k *KeyPairs) Symbol() string {
	return "ETH"
}

func (k *KeyPairs) DeriveAddress() string {
	return crypto.PubkeyToAddress(*k.publicKey).Hex()
}

func (k *KeyPairs) DerivePublicKey() string {
	return hex.EncodeToString(crypto.FromECDSAPub(k.publicKey))
}

func (k *KeyPairs) DerivePrivateKey() string {
	return hex.EncodeToString(crypto.FromECDSA(k.privateKey))
}

func (k *KeyPairs) DeriveNativeAddress() common.Address {
	return crypto.PubkeyToAddress(*k.publicKey)
}

func (k *KeyPairs) DeriveNativePrivateKey() *ecdsa.PrivateKey {
	return k.privateKey
}

func (k *KeyPairs) DeriveNativePublicKey() *ecdsa.PublicKey {
	return k.publicKey
}

// Wallet ...
type Wallet struct {
	Host     string
	client   *ethclient.Client
	Keystore []byte

	wallet *KeyPairs
}

// NewWallet ...
func NewWallet(cfg *Config) (*Wallet, error) {
	client, err := ethclient.Dial(cfg.Host)
	if err != nil {
		fmt.Println("rpc.Dial err")
		return nil, fmt.Errorf("rpc dial error: %v", err)
	}
	return &Wallet{
		Host:   cfg.Host,
		client: client,
		wallet: NewKeyPair(),
	}, nil
}

func (w *Wallet) PrintAddress() {
	fmt.Println("**************************************************")
	a0 := w.wallet.DeriveAddress()
	fmt.Printf("a0: %s\na1: %s\n", a0)
	fmt.Println("**************************************************")
}

func NewEthWalletByPath(path string, seed []byte) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, err
	}
	privKey, err := DerivePrivateKeyByPath(masterKey, path, false)
	if err != nil {
		return nil, nil, err
	}
	privateKey := privKey.ToECDSA()
	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}

func DerivePrivateKeyByPath(masterKey *hdkeychain.ExtendedKey, path string, fixIssue172 bool) (*btcec.PrivateKey, error) {
	dpath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	key := masterKey
	for _, n := range dpath {
		if fixIssue172 && key.IsAffectedByIssue172() {
			key, err = key.Derive(n)
		} else {
			key, err = key.DeriveNonStandard(n)
		}
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func DerivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}
	return publicKeyECDSA, nil
}
