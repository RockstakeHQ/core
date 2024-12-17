package mvx

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
)

func createMarketArguments() ([][]byte, error) {
	var combined []byte

	eventId := 8080
	eventIdHex := intToHex(eventId)

	description := stringToHex("FullTime Result")

	closeTimestamp := time.Now().Unix() + (2 * 60)
	closeTimestampHex := intToHex(int(closeTimestamp))

	descriptionSelection1 := intToHex(1)
	descriptionSelection2 := intToHex(2)
	descriptionSelection3 := intToHex(3)

	encoded1 := prependLength(encodeString(descriptionSelection1))
	encoded2 := prependLength(encodeString(descriptionSelection2))
	encoded3 := prependLength(encodeString(descriptionSelection3))

	combined = append(combined, encoded1...)
	combined = append(combined, encoded2...)
	combined = append(combined, encoded3...)

	eventIdBytes, err := hex.DecodeString(eventIdHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding event ID: %w", err)
	}

	descriptionBytes, err := hex.DecodeString(description)
	if err != nil {
		return nil, fmt.Errorf("error decoding description: %w", err)
	}

	closeTimestampBytes, err := hex.DecodeString(closeTimestampHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding timestamp: %w", err)
	}

	// Return arguments in the format expected by the transaction builder
	return [][]byte{
		eventIdBytes,
		descriptionBytes,
		combined,
		closeTimestampBytes,
	}, nil
}

func createMarket(
	ctx context.Context,
	txBuilder interactors.TxBuilder,
	holder core.CryptoComponentsHolder,
	eventId int,
	description string,
) error {

	smartContract := "erd1qqqqqqqqqqqqqpgqe7xe3hjzk9wt9qlz479l6pl9f69k0nuthcpsvq8r20"
	myAddress := "erd17qanf53xkluxdtvp42sjv49mwfs2emepsx047fgjvlzsjc2shcpsag4d43"

	args := blockchain.ArgsProxy{
		ProxyURL:            examples.DevnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}

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

	argsMarket, err := createMarketArguments()
	if err != nil {
		fmt.Printf("Error creating market arguments: %v\n", err)
		return err
	}

	txDataBuilder := builders.NewTxDataBuilder()
	data, err := txDataBuilder.
		Function("createMarket").
		ArgBytesList(argsMarket).
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
		GasLimit: 30000000,
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

func main() {
	ctx := context.Background()

	w := interactors.NewWallet()

	privateKey, err := w.LoadPrivateKeyFromPemFile("/Users/andrewkhirita/Desktop/rockstake/core/coverted_wallet.pem")
	if err != nil {
		panic(err)
	}

	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		panic(err)
	}

	eventId := 8000
	description := "FullTime Result"
	err = createMarket(ctx, txBuilder, holder, eventId, description)
	if err != nil {
		panic(err)
	}
}
