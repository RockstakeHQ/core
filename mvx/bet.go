package mvx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
)

func createBetArguments() ([][]byte, error) {
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
func createESDTTransferWithSCCall(
	ctx context.Context,
	txBuilder interactors.TxBuilder,
	holder core.CryptoComponentsHolder,
	tokenIdentifier string,
	amount string,
	contractAddress string,
	function string,
	args [][]byte,
) error {
	// Setup proxy
	proxyArgs := blockchain.ArgsProxy{
		ProxyURL:            examples.DevnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}

	ep, err := blockchain.NewProxy(proxyArgs)
	if err != nil {
		return fmt.Errorf("error creating proxy: %w", err)
	}

	netConfigs, err := ep.GetNetworkConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to get network configs: %w", err)
	}

	senderAddress := "erd17qanf53xkluxdtvp42sjv49mwfs2emepsx047fgjvlzsjc2shcpsag4d43" // Your sender address
	addr, err := data.NewAddressFromBech32String(senderAddress)
	if err != nil {
		return fmt.Errorf("error creating address: %w", err)
	}

	account, err := ep.GetAccount(ctx, addr)
	if err != nil {
		return fmt.Errorf("error getting account: %w", err)
	}

	// Build transaction data
	txDataBuilder := builders.NewTxDataBuilder()

	// Start with ESDTTransfer
	dataBuilder := txDataBuilder.Function("ESDTTransfer")

	// Add token identifier
	dataBuilder.ArgBytes([]byte(tokenIdentifier))

	// Add amount
	amountBig, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount format")
	}
	dataBuilder.ArgBytes(amountBig.Bytes())

	// Add function name
	dataBuilder.ArgBytes([]byte(function))

	// Add all other arguments
	for _, arg := range args {
		dataBuilder.ArgBytes(arg)
	}

	data, err := dataBuilder.ToDataBytes()
	if err != nil {
		return fmt.Errorf("error building transaction data: %w", err)
	}

	// Create transaction
	tx := &transaction.FrontendTransaction{
		Nonce:    account.Nonce,
		Value:    "0",
		Receiver: contractAddress,
		Sender:   senderAddress,
		GasPrice: netConfigs.MinGasPrice,
		GasLimit: 20000000, // Matching your Dart implementation
		Data:     data,
		ChainID:  "D", // "D" for devnet, "1" for mainnet
		Version:  1,
		Options:  0,
	}

	// Sign transaction
	err = txBuilder.ApplyUserSignature(holder, tx)
	if err != nil {
		return fmt.Errorf("error signing transaction: %w", err)
	}

	// Send transaction
	txHash, err := ep.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("error sending transaction: %w", err)
	}

	fmt.Printf("Transaction sent successfully. Hash: %s\n", txHash)
	return nil
}

func bigIntegerToHex(value *big.Int) string {
	return fmt.Sprintf("%x", value)
}

func encodeEnum(betType int) []byte {
	return []byte{byte(betType)}
}

// Example usage for placeBet
func placeBetWithESDT(
	ctx context.Context,
	txBuilder interactors.TxBuilder,
	holder core.CryptoComponentsHolder,
) error {
	// Setup your parameters
	tokenIdentifier := "MEX-a659d0"
	amount := "10000000000000000000" // 10 tokens with 18 decimals
	contractAddress := "erd17qanf53xkluxdtvp42sjv49mwfs2emepsx047fgjvlzsjc2shcpsag4d43"

	// Create bet arguments
	cid := "QmQT81GnxSyheCSfsg88efRZWgZn36Lo6By8nFEfC2pYYn"
	marketId := 15
	selectionId := 1
	odds := big.NewInt(2) // Convert your odds as needed
	betType := 0          // Back

	// Convert all arguments to bytes
	args := [][]byte{
		[]byte(cid),
		[]byte(intToHex(marketId)),
		[]byte(intToHex(selectionId)),
		[]byte(bigIntegerToHex(odds)),
		encodeEnum(betType),
		[]byte(amount),
	}

	return createESDTTransferWithSCCall(
		ctx,
		txBuilder,
		holder,
		tokenIdentifier,
		amount,
		contractAddress,
		"placeBet",
		args,
	)
}

// func main() {
// 	ctx := context.Background()

// 	// Initialize wallet
// 	w := interactors.NewWallet()
// 	privateKey, err := w.LoadPrivateKeyFromPemFile("/Users/andrewkhirita/Desktop/rockstake/core/coverted_wallet.pem")
// 	if err != nil {
// 		panic(fmt.Errorf("failed to load private key: %w", err))
// 	}

// 	// Setup crypto components
// 	suite := ed25519.NewEd25519()
// 	keyGen := signing.NewKeyGenerator(suite)

// 	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)

// 	// Create transaction builder
// 	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
// 	if err != nil {
// 		panic(fmt.Errorf("failed to create transaction builder: %w", err))
// 	}

// 	// Call the function
// 	err = placeBetWithESDT(
// 		ctx,
// 		txBuilder,
// 		holder,
// 	)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to place bet: %w", err))
// 	}
// }
