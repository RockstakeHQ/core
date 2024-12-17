package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/examples"
	"github.com/multiversx/mx-sdk-go/interactors"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

func intToHex(n int) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(n))
	return hex.EncodeToString(buf)
}

func stringToHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

func prependLength(data []byte) []byte {
	length := len(data)
	result := make([]byte, 4+length)
	binary.BigEndian.PutUint32(result[:4], uint32(length))
	copy(result[4:], data)
	return result
}

func encodeString(s string) []byte {
	return []byte(s)
}

// Test Market Creation
func createMarketArguments() ([][]byte, error) {
	var combined []byte

	eventId := 8080
	eventIdHex := intToHex(eventId)

	description := stringToHex("FullTime Result")

	closeTimestamp := time.Now().Unix() + (10 * 60)
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

// Test Bet Creation
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
		GasLimit: 30000000, // Matching your Dart implementation
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

func encodeEnum(betType int) []byte {
	return []byte{byte(betType)}
}

func hexStringToBytes(hexStr string) []byte {
	decoded, _ := hex.DecodeString(hexStr)
	return decoded
}

// Funcție helper pentru a asigura că avem întotdeauna numărul corect de cifre hex
func padHex(input string, length int) string {
	for len(input) < length {
		input = "0" + input
	}
	return input
}

func convertOddsToBigInt(odds float64) *big.Int {
	// Convert odds to BigInt following Dart logic: (odds * 100).round()
	scaled := big.NewInt(int64(odds * 100))
	return scaled
}

func bigIntegerToHex(value *big.Int) string {
	// Convert to hex string following Dart logic
	hexString := fmt.Sprintf("%x", value)
	// Ensure even length by padding with 0 if needed
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}
	return hexString
}

const (
	Back = 0
	Lay  = 1
)

func placeBetWithESDT(
	ctx context.Context,
	txBuilder interactors.TxBuilder,
	holder core.CryptoComponentsHolder,
) error {
	tokenIdentifier := "MEX-a659d0"
	betAmount := "20000000000000000000" // 10 tokens with 18 decimals
	contractAddress := "erd1qqqqqqqqqqqqqpgqe7xe3hjzk9wt9qlz479l6pl9f69k0nuthcpsvq8r20"

	// Convert CID to hex
	cid := "QmQT81GnxSyheCSfsg88efRZWgZn36Lo6By8nFEfC2pYYn"
	cidHex := hex.EncodeToString([]byte(cid))

	// Convert market ID to hex (15)
	marketId := 5
	marketIdHex := intToHex(marketId)

	// Convert selection ID to hex (1)
	selectionId := 1
	selectionIdHex := intToHex(selectionId)

	// Convert odds following Dart logic
	odds := 2.0
	oddsForContract := convertOddsToBigInt(odds)
	oddsHex := bigIntegerToHex(oddsForContract)

	betType := Lay
	var liability *big.Int

	if betType == Back {
		liability = big.NewInt(0)
	} else {
		// Pentru Lay, calculăm liability
		liability = big.NewInt(10)
	}

	liabilityHex := bigIntegerToHex(liability)
	betTypeBytes := encodeEnum(betType)

	// Prepare arguments
	args := [][]byte{
		hexToBytes(cidHex),
		hexToBytes(marketIdHex),
		hexToBytes(selectionIdHex),
		hexToBytes(oddsHex),
		betTypeBytes,
		hexToBytes(liabilityHex), // Folosim liability în loc de amount pentru ultimul argument
	}

	return createESDTTransferWithSCCall(
		ctx,
		txBuilder,
		holder,
		tokenIdentifier,
		betAmount,
		contractAddress,
		"placeBet",
		args,
	)
}

// Helper function to convert hex string to bytes
func hexToBytes(hexStr string) []byte {
	decoded, _ := hex.DecodeString(hexStr)
	return decoded
}

// Main for Creating Market
// func main() {
// 	ctx := context.Background()

// 	w := interactors.NewWallet()

// 	privateKey, err := w.LoadPrivateKeyFromPemFile("/Users/andrewkhirita/Desktop/rockstake/core/converted_wallet.pem")
// 	if err != nil {
// 		panic(err)
// 	}

// 	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)
// 	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
// 	if err != nil {
// 		panic(err)
// 	}

// 	eventId := 8000
// 	description := "FullTime Result"
// 	err = createMarket(ctx, txBuilder, holder, eventId, description)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// Main for Creating Bet
func main() {
	ctx := context.Background()

	// Initialize wallet
	w := interactors.NewWallet()
	privateKey, err := w.LoadPrivateKeyFromPemFile("/Users/andrewkhirita/Desktop/rockstake/core/converted_wallet.pem")
	if err != nil {
		panic(fmt.Errorf("failed to load private key: %w", err))
	}

	// Setup crypto components
	suite := ed25519.NewEd25519()
	keyGen := signing.NewKeyGenerator(suite)

	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, privateKey)

	// Create transaction builder
	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		panic(fmt.Errorf("failed to create transaction builder: %w", err))
	}

	// Call the function
	err = placeBetWithESDT(
		ctx,
		txBuilder,
		holder,
	)
	if err != nil {
		panic(fmt.Errorf("failed to place bet: %w", err))
	}
}
