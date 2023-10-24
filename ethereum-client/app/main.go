package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)
// remix 部署合约： https://remix.ethereum.org/#lang=en&optimize=false&runs=200&evmVersion=null&version=soljson-v0.7.0+commit.9e61f92b.js
// 广播交易查询结果： https://goerli.etherscan.io/tx/0xdf9d6b7935b52b7766dcc22493bdd2df8619a03d7805cabd08e218673a8bb399
// eth decode tx：  https://flightwallet.github.io/decode-eth-tx/
// 水龙头： https://blog.csdn.net/cljdsc/article/details/130641872
// https://faucet.quicknode.com/ethereum/goerli/?transactionHash=0x4487e8214c2616f3a41fc63838d0e219483471d0c182ca9dcc8013011fb783b7 水龙头
// 水龙头领取地址： https://goerlifaucet.com/
func main() {
	//mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	//wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/100001")
	////path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	//account, err := wallet.Derive(path, true)
	//if err != nil {
	//	log.Fatal(err)
	//}

	client, err := ethclient.Dial("https://rpc.ankr.com/eth_goerli/8b4a7aff54ac22cd3d15d0e58b3ba1a6ee3f90b2233cba73bd7093dbcfe885dd")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("7efa92903c5e51228cb7147125c483efd82d1fd6af06398da09ea2cbf9fb9fff")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress("0xbc6Bb41Cf8071BCAE219097fFD0159F75d609d37"))
	if err != nil {
		log.Fatal(err)
	}


	//nonce := uint64(0)
	value := big.NewInt(0)
	gasLimit := uint64(52141)
	gasPrice := big.NewInt(21000000000)

	//800000000000000000000000
	toAddress := common.HexToAddress("0x7956Cae0463572955c032AE2CF857fCb9f3D7c9c")
	tokenAddress := common.HexToAddress("0x49aa3681d1ce3a87ee675d533c20ca92a4262e15")

	fmt.Println("fromAddress: "+ fromAddress.String())
	fmt.Println(nonce)
	fmt.Println("toAddress: "+ toAddress.String())
	fmt.Println("tokenAddress: "+ tokenAddress.String())

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	amount := new(big.Int)
	amount.SetString("1000000000000000000", 10) // sets the value to 1000 tokens, in the token denomination
	//chainID := big.NewInt(5)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)


	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)



	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	data, err = signedTx.MarshalBinary()
	fmt.Println("待广播str : "+hexutil.Encode(data))
	fmt.Println("Hash :" + signedTx.Hash().String())
	fmt.Println(signedTx)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("区块浏览器：https://goerli.etherscan.io/tx/" + signedTx.Hash().String())
	spew.Dump(signedTx)
}


