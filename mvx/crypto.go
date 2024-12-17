package mvx

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"

	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
)

// var (
// 	suite  = ed25519.NewEd25519()
// 	keyGen = signing.NewKeyGenerator(suite)
// )

func burnTokens(
	ctx context.Context,
	txBuilder interactors.TxBuilder,
	holder core.CryptoComponentsHolder,
	tokenIdentifier string,
	amount *big.Int,
) error {
	// Setup proxy
	args := blockchain.ArgsProxy{
		ProxyURL:            examples.DevnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}

	smartContract := os.Getenv("SMART_CONTRACT")
	myAddress := os.Getenv("MY_ADDRESS")

	ep, err := blockchain.NewProxy(args)
	if err != nil {
		return fmt.Errorf("error creating proxy: %w", err)
	}

	netConfigs, err := ep.GetNetworkConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to get network configs: %w", err)
	}

	addr, err := data.NewAddressFromBech32String(myAddress)
	if err != nil {
		return fmt.Errorf("error creating address: %w", err)
	}

	account, err := ep.GetAccount(ctx, addr)
	if err != nil {
		return fmt.Errorf("error getting account: %w", err)
	}

	txDataBuilder := builders.NewTxDataBuilder()
	data, err := txDataBuilder.
		Function("burnTokens").
		ArgBytes([]byte(tokenIdentifier)).
		ArgInt64(0).
		ArgBigInt(amount).
		ToDataBytes()
	if err != nil {
		return fmt.Errorf("error building transaction data: %w", err)
	}

	// Construim tranzacția
	tx := &transaction.FrontendTransaction{
		Nonce:    account.Nonce,
		Value:    "0",
		Receiver: smartContract,
		Sender:   myAddress,
		GasPrice: netConfigs.MinGasPrice,
		GasLimit: 5000000,
		Data:     data,
		ChainID:  "D",
		Version:  1,
		Options:  0,
	}

	// Aplicăm semnătura folosind txBuilder
	err = txBuilder.ApplyUserSignature(holder, tx)
	if err != nil {
		return fmt.Errorf("error signing transaction: %w", err)
	}

	// Trimitem tranzacția
	txHash, err := ep.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("error sending transaction: %w", err)
	}

	fmt.Printf("Transaction sent successfully. Hash: %s\n", txHash)
	return nil
}

// func main() {
// 	ctx := context.Background()

// 	w := interactors.NewWallet()

// 	privateKey, err := w.LoadPrivateKeyFromPemFile("/Users/andrewkhirita/Desktop/rockstake/core/wallets/wallet_shard2_2.pem")
// 	if err != nil {
// 		panic(err)
// 	}

// 	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
// 	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
// 	if err != nil {
// 		panic(err)
// 	}

// 	tokenId := "SNOW-d7a8f5"
// 	// Cu aceasta
// 	amount := new(big.Int)
// 	amount.SetString("5000000000000000", 10) // Folosim string pentru numere mari

// 	err = burnTokens(ctx, txBuilder, holder, tokenId, amount)
// 	if err != nil {
// 		panic(err)
// 	}
// }
